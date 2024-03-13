package consts


const (
	ExchangeType        = "topic"
	NOVA                = "nova"
	INSTANCE            = "instance"
	IRONIC              = "ironic"
	IRONICFLAVOR        = "SEPC"
	InstanceId          = "instance_id"

	CINDER              = "cinder"
	VOLUME              = "volume"
	SNAPSHOT     		= "snapshot"
	BACKUP              = "backup"
	VolumeId            = "volume_id"
	SnapshotId          = "snapshot_id"
	BackupId            = "backup_id"

	NEUTRON             = "neutron"
	FIREWALLGROUP       = "firewall_group"
    PORT                = "port"
    FLOATINGIPPORT      = "floatingip_port"
    FLOATINGIP          = "floatingip"
    FIREWALL            = "firewall"

	COMPUTE             = "compute"
	BAREMETALNODE       = "baremetal_node"
	ComputeService      = "compute_service"
	Binary              = "binary"
	NovaCompute         = "nova-compute"

	Id                  = "id"
	Uuid                = "uuid"
	State               = "state"
	Status              = "status"
	PowerState          = "power_state"
	ProvisionState      = "provision_state"
	Maintenance         = "maintenance"
	Deleted             = "deleted"
	Available           = "available"
	InUse               = "in-use"
	Enabled             = "enabled"
	Disabled            = "disabled"
    Active              = "ACTIVE"
    Down                = "DOWN"
    DeviceOwner         = "device_owner"
    FipDeviceOwner      = "network:floatingip"

	FixedIps            = "fixed_ips"
	Host                = "host"
	Node                = "node"
	VolumeAttachment    = "volume_attachment"
	Size                = "size"
	Exception           = "exception"
	Reason              = "reason"
	DisabledReason      = "disabled_reason"
	Error               = "error"
	Stopped             = "stopped"
	ErrorMsg            = "error_msg"
	InstanceType        = "instance_type"

	NovaSuffix          = "/inner/server/updateResourceState"
	CinderSuffix        = "/inner/integration/ebs/stateMachine/pushStateMessage"
	NeutronSuffix       = "/inner/network/updateResourceState"
	ComputeSuffix       = "/inner/server/updateNodeResourceState"
	GetMqInfoSuffix     = "/inner/regions/mqinfos"

	SERVICEACCOUNT      = "/var/run/secrets/kubernetes.io/serviceaccount"
	NamespacePath       = SERVICEACCOUNT + "/namespace"
)

var EventStates = map[string]string{
	"compute.instance.create.end": NOVA,
	"compute.instance.create.error": NOVA,
	"compute.instance.delete.end": NOVA,
	"compute.instance.soft_delete.end": NOVA,
	"compute.instance.power_off.end": NOVA,
	"compute.instance.power_on.end": NOVA,
	"compute.instance.reboot.end": NOVA,
	"compute.instance.reboot.error": NOVA,
	"compute.instance.rebuild.end": NOVA,
	"compute.instance.rebuild.error": NOVA,
	"compute.instance.resize.error": NOVA,
	"compute.instance.resize.confirm.end": NOVA,
	"compute.instance.live_migration._post.end": NOVA,
	"compute_task.build_instances": NOVA,
	"compute_task.migrate_server": NOVA,
	"compute_task.rebuild_server": NOVA,
	"finish_resize": NOVA,
	"compute.instance.live_migration.rollback.dest.end": NOVA,
	// live resize
	"instance.live_resize.end": NOVA,
	"compute_task.live_resize": NOVA,

	"baremetal.node.create.end": BAREMETALNODE,
	"baremetal.node.create.error": BAREMETALNODE,
	"baremetal.node.power_set.end": BAREMETALNODE,
	"baremetal.node.power_state_corrected.success": BAREMETALNODE,
	"baremetal.node.power_set.error": BAREMETALNODE,
	"baremetal.node.delete.end": BAREMETALNODE,
	"baremetal.node.delete.error": BAREMETALNODE,
	"baremetal.node.provision_set.end": BAREMETALNODE,
	"baremetal.node.provision_set.success": BAREMETALNODE,
	"baremetal.node.provision_set.error": BAREMETALNODE,
	"baremetal.node.update.end": BAREMETALNODE,
	"baremetal.node.update.error": BAREMETALNODE,
	"baremetal.node.maintenance_set.end": BAREMETALNODE,
	"baremetal.node.maintenance_set.error": BAREMETALNODE,
	"service.create": ComputeService,
	"service.update": ComputeService,
	"service.delete": ComputeService,

	"scheduler.create_volume": VOLUME,         //create failed
	"volume.create.end": VOLUME,
	"volume.delete.end": VOLUME,
	"volume.delete.error": VOLUME,
	"volume.attach.end": VOLUME,
	"volume.detach.end": VOLUME,
	"volume.retype": VOLUME,
	"volume.resize.end": VOLUME,
	"volume.resize.error": VOLUME,
	"compute.instance.volume.detach": VOLUME,
	"attach_volume": VOLUME,
	"detach_volume": VOLUME,
	"snapshot.create.end": SNAPSHOT,
	"snapshot.create.error": SNAPSHOT,
	"snapshot.revert.end": SNAPSHOT,
	"volume.revert.end": VOLUME,
	"snapshot.delete.end": SNAPSHOT,
	"snapshot.delete.error": SNAPSHOT,

	"firewall_group.create.end": FIREWALLGROUP,
	"firewall_group.update.end": FIREWALLGROUP,
	"firewall_group.delete.end": FIREWALLGROUP,
	"firewall_group.update_status": FIREWALLGROUP,
    "port.update.end": PORT,
    "floatingip.create.end": FLOATINGIP,
    "floatingip.update.end": FLOATINGIP,
    "floatingip.delete.end": FLOATINGIP,
    "sdn.floatingip.update.end": FLOATINGIP,
    "sdn.firewall.update.end": FIREWALL,
}

var SecondaryMap = map[string]string{
	VOLUME: CINDER,
	SNAPSHOT: CINDER,
	BACKUP: CINDER,
	FIREWALLGROUP: NEUTRON,
	PORT: NEUTRON,
	FLOATINGIP: NEUTRON,
	FIREWALL: NEUTRON,
	IRONIC: NOVA,
	INSTANCE: NOVA,
	BAREMETALNODE: COMPUTE,
	ComputeService: COMPUTE,
}

var ResourceIdMap = map[string]string{
	NOVA: InstanceId,
	VOLUME: VolumeId,
	SNAPSHOT: SnapshotId,
	BACKUP: BackupId,
	FIREWALLGROUP: Id,
	PORT: Id,
	FLOATINGIP: Id,
	FIREWALL: Id,
}

var ResourceStateMap = map[string]string{
	NOVA: State,
	VOLUME: Status,
	SNAPSHOT: Status,
	BACKUP: Status,
	FIREWALLGROUP: Status,
	PORT: Status,
	FLOATINGIP: Status,
	FIREWALL: Status,
}
