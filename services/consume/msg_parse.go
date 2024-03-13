package consume

import (
	"go.uber.org/zap"
	"log"
	"regexp"
	"sincerecloud.com/openstack_notifier/consts"
	"sincerecloud.com/openstack_notifier/global"
	"sincerecloud.com/openstack_notifier/services/publish"
	"strings"
)

// InstanceCreateErrorMsgFormatters For create.error event, these are BuildAbortException message
var InstanceCreateErrorMsgFormatters = []string{
	"Build of instance .* aborted",
	"Maximum number of fixed IPs exceeded",
	"No more available networks",
	"No fixed IP addresses available for network",
	"Virtual Interface creation failed",
	"Creation of virtual interface with unique mac address failed",
	"The fixed IP associated with port .* is not compatible with the host",
	"Unable to automatically allocate a network for project",
	"Using networks with QoS policy is not supported for instance",
	"The created instance's disk would be too small",
	"Flavor's memory is too small for requested image",
	"Image .* is not active",
	"Image .* is unacceptable",
	"Disk info file is invalid",
	"Disk format .* is not acceptable",
	"Image signature certificate validation failed for certificate",
	"Volume encryption is not supported for .* volume",
	"Invalid input received",
	"The requested amount of video memory .* is higher than the maximum allowed by flavor",
	"Signature verification for the image failed",
}


type Parser struct {
	eventType           string
	resource            string
	key                 string
	msgInfo             map[string]interface{}
	publishMsg          *publish.PublishMsg
	IsSDN               bool
}

func NewParser(eventType, resource, key string, msgInfo map[string]interface{}, isSdn bool) *Parser {
	return &Parser{
		eventType: eventType,
		resource: resource,
		key: key,
		msgInfo: msgInfo,
		publishMsg: new(publish.PublishMsg),
		IsSDN: isSdn,
	}
}

func (p *Parser) Parse() *publish.PublishMsg {
	switch p.key {
	case consts.NOVA:
		return p.handleNova()
	case consts.CINDER:
		return p.handleCinder()
	case consts.NEUTRON:
		return p.handleNeutron()
	case consts.COMPUTE:
		return p.handleCompute()
	default:
		return nil
	}
}

func (p *Parser) isMatchMsg(message string) bool {
	for _, s := range InstanceCreateErrorMsgFormatters {
		isMatch, _ := regexp.MatchString(s, message)
		if isMatch {
			return true
		}
	}
	return false
}

func (p *Parser) handleNova() *publish.PublishMsg {
	idName := consts.ResourceIdMap[p.resource]
	stateName := consts.ResourceStateMap[p.resource]
	payload := p.msgInfo["payload"].(map[string]interface{})

	var flavorId string
	var flavorName string
	var resourceId string
	var resourceState string
	var errorMsg interface{}
	var host string
	var node string
	if p.eventType == "finish_resize" {   // finish_resize
		resourceId = p.traceStringVal(payload, []string{"args", "instance", "uuid"})
		flavorName = p.traceStringVal(payload, []string{"args", "instance", "flavor", "name"})
		flavorId = p.traceStringVal(payload, []string{"args", "instance", "flavor", "flavorid"})
		resourceState = p.traceStringVal(payload, []string{"args", "instance", "vm_state"})
		host = p.traceStringVal(payload, []string{"args", consts.Host})
		node = p.traceStringVal(payload, []string{"args", consts.Node})
		errorMsg = p.traceStringVal(payload, []string{consts.Exception})
	} else if p.eventType == "compute_task.build_instances" ||
		p.eventType == "compute_task.migrate_server" ||
		p.eventType == "compute_task.rebuild_server" ||
		p.eventType == "compute_task.live_resize" {
		log.Println(payload)
		log.Println(idName)
		resourceId = payload[idName].(string)
		flavorName = p.traceStringVal(payload, []string{"request_spec", "instance_type", "name"})
		flavorId = p.traceStringVal(payload, []string{"request_spec", "instance_type", "flavorid"})
		resourceState = payload[stateName].(string)
		// no host
		// no node
		errorMsg = p.traceStringVal(payload, []string{consts.Reason})
	} else if p.eventType == "instance.live_resize.end" {   // finish_resize
		resourceId = p.traceStringVal(payload, []string{"nova_object.data", "uuid"})
		flavorName = p.traceStringVal(payload, []string{"nova_object.data", "flavor", "nova_object.data", "name"})
		flavorId = p.traceStringVal(payload, []string{"nova_object.data", "flavor", "nova_object.data", "flavorid"})
		resourceState = p.traceStringVal(payload, []string{"nova_object.data", "instance", "state"})
		host = p.traceStringVal(payload, []string{"nova_object.data", consts.Host})
		node = p.traceStringVal(payload, []string{"nova_object.data", consts.Node})
	} else {
		resourceId = payload[idName].(string)
		flavorName = p.traceStringVal(payload, []string{consts.InstanceType})
		flavorId = p.traceStringVal(payload, []string{"instance_flavor_id"})
		host = p.traceStringVal(payload, []string{consts.Host})
		node = p.traceStringVal(payload, []string{consts.Node})
		errorMsg, _ = payload[consts.Exception]
		if errorMsg != nil {
			resourceState = consts.Error
			if p.eventType == "compute.instance.create.error" && !p.isMatchMsg(errorMsg.(string)) {
				return nil       // match BuildAbortException, send create.error msg
			} else if p.eventType == "compute.instance.resize.error" {
				errorMsgStr := errorMsg.(string)
				if payload[stateName].(string) != consts.Error && !strings.Contains(errorMsgStr, "NoValidHost_Remote") {
					return nil
				}
				resourceState = payload[stateName].(string)
			}
		} else {
			resourceState = payload[stateName].(string)
		}

	}

	//ResourceId
	p.publishMsg.ResourceId = resourceId
	//ResourceType ironic/nova
	p.publishMsg.ResourceType = consts.INSTANCE
	if isIronic(flavorName) {
		p.publishMsg.ResourceType = consts.IRONIC
	}
	//ResourceState
	p.publishMsg.ResourceState = resourceState
	// FixedIp
	if fixedIps, exist := payload[consts.FixedIps]; exist {
		var resFixedIps []string
		for _, fixedIp := range fixedIps.([]interface{}) {
			if address, exist := fixedIp.(map[string]interface{})["address"]; exist {
				resFixedIps = append(resFixedIps, address.(string))
			}
		}
		p.publishMsg.FixedIps = resFixedIps
	}
	//Host
	p.publishMsg.Host = host
	//Node
	p.publishMsg.Node = node
	//ErrorMsg
	p.publishMsg.ErrorMsg = errorMsg
	//FlavorId
	p.publishMsg.FlavorId = flavorId
	return p.publishMsg
}

func (p *Parser) handleCinder() *publish.PublishMsg {
	payload := p.msgInfo["payload"].(map[string]interface{})

	var resourceId string
	var resourceState string
	idName := consts.ResourceIdMap[p.resource]
	stateName := consts.ResourceStateMap[p.resource]
	if p.eventType == "attach_volume" {
		resourceId = p.traceStringVal(payload, []string{"args", "bdm", consts.VolumeId})
		resourceState = consts.Available
		p.publishMsg.ErrorMsg = p.traceStringVal(payload, []string{consts.Exception})
	} else if p.eventType == "detach_volume" {
		resourceId = p.traceStringVal(payload, []string{"args", consts.VolumeId})
		resourceState = consts.InUse
		p.publishMsg.ErrorMsg = p.traceStringVal(payload, []string{consts.Exception})
	} else if p.eventType == "scheduler.create_volume" {
		resourceId = payload[idName].(string)
		resourceState = payload[consts.State].(string)
        p.publishMsg.ErrorMsg = p.traceStringVal(payload, []string{consts.Reason})
	} else if p.eventType == "compute.instance.volume.detach" {
		resourceId = payload[idName].(string)
		resourceState = consts.Available
	} else {
		resourceId = payload[idName].(string)
		resourceState = payload[stateName].(string)
		if p.eventType == "volume.detach.end" {
			attachment := payload[consts.VolumeAttachment].([]interface{})
			if len(attachment) == 0 {
				resourceState = consts.Available
			}
		}
		if errorMsg, exist := payload[consts.ErrorMsg]; exist {
			p.publishMsg.ErrorMsg = errorMsg
		}
    }
    // ResourceId
    p.publishMsg.ResourceId = resourceId

    // ResourceState
    p.publishMsg.ResourceState = resourceState

	// ResourceType
	p.publishMsg.ResourceType = p.resource

	// VolumeAttachment
	if volumeAttachment, exist := payload[consts.VolumeAttachment]; exist {
		var resVolumeAttachments []string
		for _, attachment := range volumeAttachment.([]interface{}) {
			if instance, exist := attachment.(map[string]interface{})["instance_uuid"]; exist {
				resVolumeAttachments = append(resVolumeAttachments, instance.(string))
			}
		}
		p.publishMsg.VolumeAttachment = resVolumeAttachments
	}

	//VolumeSize
	if p.eventType == "volume.resize.end" {
		if size, exist := payload[consts.Size]; exist {
			p.publishMsg.VolumeSize = size.(float64)
		}
	}
	return p.publishMsg
}


func (p *Parser) handleNeutron() *publish.PublishMsg {
	idName := consts.ResourceIdMap[p.resource]
	stateName := consts.ResourceStateMap[p.resource]

	temp := p.msgInfo["payload"].(map[string]interface{})
	v, ok := temp[p.resource]
	if !ok {
		global.LOG.Error("Failed to parse payload for", zap.String("resource", p.resource))
		return nil
	}
	payload := v.(map[string]interface{})

	resourceType := p.resource
	//ResourceState
	var state string
	state = payload[stateName].(string)
	if strings.Contains(p.eventType, "delete.end") {
		state = consts.Deleted
	}
	switch p.eventType {
	case "port.update.end":
		// For port, only network:floatingip port, filter port.update.end in non-snd env
		if payload[consts.DeviceOwner].(string) != consts.FipDeviceOwner || !p.IsSDN {
			return nil
		} else {
			resourceType = consts.FLOATINGIPPORT
		}
	case "floatingip.update.end":
		// filter floatingip.update.end in sdn env
		if p.IsSDN {
			return nil
		} else {
			if portId, ok := payload["port_id"]; ok {
				if portId == nil {    // The state is DOWN where fip disassociates port
					state = consts.Down
				} else {  // The state is ACTIVE where fip associates port
					state = consts.Active
				}
			}
		}
	}

	p.publishMsg.ResourceState = state

	//ErrorMsg
	if errorMsg, exist := payload[consts.ErrorMsg]; exist {
		if errorMsg != nil {
			p.publishMsg.ErrorMsg = errorMsg.(interface{})
		}
	}

	p.publishMsg.ResourceId = payload[idName].(string)
	p.publishMsg.ResourceType = resourceType
	return p.publishMsg
}

func (p *Parser) handleCompute() *publish.PublishMsg {
	payload := p.msgInfo["payload"].(map[string]interface{})

	var resourceId string
	var resourceState string
	var maintenance bool
	var provisionState string
	var resourceType string
	var errorMsg interface{}
	var host string
	if strings.HasPrefix(p.eventType, "baremetal") {
		resourceId = p.traceStringVal(payload, []string{"ironic_object.data", consts.Uuid})
		resourceState = p.traceStringVal(payload, []string{"ironic_object.data", consts.PowerState})
		provisionState = p.traceStringVal(payload, []string{"ironic_object.data", consts.ProvisionState})
		resourceType = consts.BAREMETALNODE
		errorMsg = p.traceStringVal(payload, []string{"ironic_object.data", "last_error"})
		maintenance = p.traceBoolVal(payload, []string{"ironic_object.data", consts.Maintenance})
	} else if strings.HasPrefix(p.eventType, "service") {
		resourceId = p.traceStringVal(payload, []string{"nova_object.data", consts.Uuid})
		binary := p.traceStringVal(payload, []string{"nova_object.data", consts.Binary})
		if binary != consts.NovaCompute {
			return nil
		}
		disabled := p.traceBoolVal(payload, []string{"nova_object.data", consts.Disabled})
		if disabled == true {
			provisionState = consts.Disabled
		} else {
			provisionState = consts.Enabled
		}
		resourceType = consts.ComputeService
		host = p.traceStringVal(payload, []string{"nova_object.data", consts.Host})
		errorMsg = p.traceStringVal(payload, []string{"nova_object.data", consts.DisabledReason})
	}
	p.publishMsg.ResourceId = resourceId
	p.publishMsg.ResourceType = resourceType
	p.publishMsg.ResourceState = resourceState
	p.publishMsg.ErrorMsg = errorMsg
	p.publishMsg.ProvisionState = provisionState
	p.publishMsg.Maintenance = maintenance
	p.publishMsg.Host = host
	return p.publishMsg
}

func (p *Parser) traceStringVal(payload map[string]interface{}, path []string) string {
	var value interface{}
	var exist bool
	for _, key := range path {
		if value, exist = payload[key]; exist {
			switch value.(type) {
			case string:
				return value.(string)
			case map[string]interface{}:
				payload = value.(map[string]interface{})
			}
		}
	}
	return ""
}

func (p *Parser) traceBoolVal(payload map[string]interface{}, path []string) bool {
	var value interface{}
	var exist bool
	for _, key := range path {
		if value, exist = payload[key]; exist {
			switch value.(type) {
			case bool:
				return value.(bool)
			case map[string]interface{}:
				payload = value.(map[string]interface{})
			}
		}
	}
	return false
}

func isIronic(flavorName string) bool {
	if strings.HasPrefix(flavorName, consts.IRONICFLAVOR) {
		return true
	}
	return false
}
