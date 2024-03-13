package consume

import (
	"fmt"
	"reflect"
	"sincerecloud.com/openstack_notifier/consts"
	"testing"
)


func TestHandleNova_creat_end(t *testing.T) {
	jsonBody := []byte(`{"payload": {"state_description": "","availability_zone": "nova", "terminated_at": "", "ephemeral_gb": 0, "instance_type_id": 55,"message": "Success", "deleted_at": "", "fixed_ips": [{"version": 4,"vif_mac": "fa:16:3e:ab:77:b8", "floating_ips": [], "label": "rds-mgnt-vlan529","meta": {}, "address": "10.50.29.99", "type": "fixed"}],"instance_id": "640ee2c0-c183-4a24-8345-9cfb712408a9", "display_name": "dx_instance1","reservation_id": "r-sqvxf2sq", "hostname": "dx-instance1", "state": "active","progress": "", "launched_at": "2023-01-19T02:48:55.384192", "metadata": {},"node": "dev-rds.vim1.local", "ramdisk_id": "", "access_ip_v6": null, "disk_gb": 0,"access_ip_v4": null, "kernel_id": "", "host": "dev-rds.vim1.local","user_id": "83bac5ab80b14f74b06f157daa6fd64b","image_ref_url": "http://10.50.1.57:9292/images/032e83c9-eddb-40d8-8c23-e4e8015c79ee", "cell_name": "","root_gb": 0, "tenant_id": "a01edbe369764a6ba25798bb477f245b", "created_at": "2023-01-19 02:48:46+00:00","memory_mb": 1024, "instance_type": "SEPC_flavor","vcpus": 2, "image_meta": {"hw_qemu_guest_agent": "yes", "os_distro": "centos","image_type": "image", "container_format": "bare", "min_ram": "0","owner_specified.openstack.sha256": "", "disk_format": "raw", "os_admin_user": "root","usage_type": "common", "owner_specified.openstack.object": "images/centos7.9-220929","owner_specified.openstack.md5": "", "min_disk": "0", "os_type": "linux","base_image_ref": "032e83c9-eddb-40d8-8c23-e4e8015c79ee"}, "architecture": null, "os_type": "linux","instance_flavor_id": "02ed66eb-9002-4336-ae23-e5ff1efad435"}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("compute.instance.create.end", "nova", "nova", msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != "ironic" || publisher.ErrorMsg != nil ||
		!reflect.DeepEqual(publisher.FixedIps, []string{"10.50.29.99"}) ||
		publisher.Host != "dev-rds.vim1.local" || publisher.FlavorId != "02ed66eb-9002-4336-ae23-e5ff1efad435"{
		t.Fatal("Failed to handleNova")
	}
}

func TestHandleNova_creat_error(t *testing.T) {
	jsonBody := []byte(`{"payload": {"state_description": "spawning", "code": 500, "availability_zone": "nova", "terminated_at": "", "ephemeral_gb": 0, "instance_type_id": 1069,"message": "Virtual Interface creation failed", "deleted_at": "", "reservation_id": "r-gs4q4g6o","instance_id": "ccca82ab-febc-4647-9ed5-e1081d186e43", "display_name": "volume_type","hostname": "volume-type", "state": "building", "progress": "", "launched_at": "","metadata": {}, "node": "con01.vim1.local", "ramdisk_id": "", "access_ip_v6": null, "disk_gb": 0,"access_ip_v4": null, "kernel_id": "", "host": "con01.vim1.local","user_id": "a1e0d70fc2274d39af8434d633c3347e", "image_ref_url": "http://10.50.31.1:9292/images/","cell_name": "", "exception": "{'message': u'Virtual Interface creation failed', 'class': 'VirtualInterfaceCreateException','kwargs': {'code': 500}}", "root_gb": 0, "tenant_id": "7e8babd4464e4c6da382a1a29d8da53a","created_at": "2023-02-07 07:25:42+00:00", "memory_mb": 8192, "instance_type": "004008","vcpus": 4, "image_meta": {"os_distro": "centos", "image_type": "image","container_format": "bare", "hw_qemu_guest_agent": "yes", "disk_format": "raw","os_admin_user": "root", "os_version": "7.9", "min_ram": "1024", "min_disk": "8","usage_type": "common", "base_image_ref": ""}, "architecture": null, "os_type": null,"instance_flavor_id": "004008"}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("compute.instance.create.error", "nova","nova", msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != "instance" || publisher.ErrorMsg == nil ||
		publisher.ResourceState != "error" || publisher.Host != "con01.vim1.local" ||
		publisher.FlavorId != "004008"{
		t.Fatal("Failed to handleNova")
	}
}

func TestHandleNova_finish_resize(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"exception": "libvirtError(\"internal error: qemu unexpectedly closed the monitor: 202-02-27T06:55:58.362653Z qemu-kvm: -drive file=rbd:volumes/volume-323c67c7-6cca-4c30-af40-48b14d5221c2:id=cinder:auth_supported=cephx;none:mon_host=10.50.31.1:6789,file.password-secret=virtio-disk0-secret0,format=raw,if=none,id=drive-virtio-disk0,serial=323c67c7-6cca-4c30-af40-48b14d5221c2,cache=writeback,discard=unmap: 'serial' is deprecated, please use the corresponding option of '-device' insteadn2023-02-27T06:55:58.412089Z qemu-kvm: cannot set up guest memory 'pc.ram': Cannot allocate memory\",)", "args": {"instance": {"vm_state": "error", "pci_requests": {"instance_uuid": "a55c993a-4012-4b41-97d3-69d0666e8197", "requests": []}, "availability_zone": "nova", "terminated_at": null, "ephemeral_gb": 0, "old_flavor": {"memory_mb": 4096, "root_gb": 0, "deleted_at": null, "description": null, "deleted": false, "created_at": "2022-04-01T09:14:22.000000", "ephemeral_gb": 0, "updated_at": null, "disabled": false, "vcpus": 1, "extra_specs": {":category": "memory_optimized", "hw:live_resize": "True", ":architecture": "x86_architecture", "hw:numa_nodes": "1"}, "swap": 0, "rxtx_factor": 1.0, "is_public": true, "flavorid": "0348236c-417d-4409-87db-36550fdeebe8", "vcpu_weight": 0, "id": 801, "name": "C1G4"}, "updated_at": "2023-02-27T06:55:59.000000", "numa_topology": null, "cleaned": true, "vm_mode": null, "flavor": {"memory_mb": 4096, "root_gb": 0, "deleted_at": null, "description": null, "deleted": false, "created_at": "2022-04-01T09:14:22.000000", "ephemeral_gb": 0, "updated_at": null, "disabled": false, "vcpus": 1, "extra_specs": {":category": "memory_optimized", "hw:live_resize": "True", ":architecture": "x86_architecture", "hw:numa_nodes": "1"}, "swap": 0, "rxtx_factor": 1.0, "is_public": true, "flavorid": "0348236c-417d-4409-87db-36550fdeebe8", "vcpu_weight": 0, "id": 801, "name": "SEPC_C1G4"}, "deleted_at": null, "reservation_id": "r-qihzjdy2", "id": 11696, "security_groups": [], "disable_terminate": false,"user_id": "a1e0d70fc2274d39af8434d633c3347e", "uuid": "a55c993a-4012-4b41-97d3-69d0666e8197", "default_swap_device": null}, "info_cache": {"_obj_instance_uuid": "a55c993a-4012-4b41-97d3-69d0666e8197", "_changed_fields": [], "_obj_updated_at": "2023-02-27T04:01:54.000000", "VERSION": "1.5", "_obj_network_info": [{"profile": {}, "ovs_interfaceid": "269fef65-7a26-427e-8249-7f3dbea71c29", "preserve_on_delete": false, "network": {"bridge": "br-int", "label": "dx_net1", "meta": {"injected": false, "tunneled": true, "tenant_id": "7e8babd4464e4c6da382a1a29d8da53a", "physical_network": null, "mtu": 1450}, "id": "f84b3aeb-f8c8-489c-a8ad-8d14b011db58", "subnets": [{"ips": [{"meta": {}, "type": "fixed", "version": 4, "address": "23.4.1.129", "floating_ips": []}], "version": 4, "meta": {"dhcp_server": "23.4.1.2"}, "dns": [], "routes": [], "cidr": "23.4.1.0/24", "gateway": {"meta": {}, "type": "gateway", "version": 4, "address": "23.4.1.1"}}]}, "devname": "tap269fef65-7a", "qbh_params": null, "vnic_type": "normal", "meta": {}, "details": {"ovs_hybrid_plug": false, "bridge_name": "br-int", "datapath_type": "system", "port_filter": true, "connectivity": "l2"}, "address": "fa:16:3e:e5:a9:2e", "active": true, "type": "ovs", "id": "269fef65-7a26-427e-8249-7f3dbea71c29", "qbg_params": null}], "_context": {"service_user_domain_name": null, "service_user_id": null, "auth_token": "gAAAAABj_CrUITIzmjumgkA3iaAICzaDqNNHxfcaA3pF6KQzrOjJjPRsHmy-9UJ52pu2Befpe0GqNfOZpU4GLHJELLTGai8e8YxhHdK-cN2tA2XAU_WhmgO5z8TPrKfeZBPzbo3jRrCwDzoOcE_piPmdu39dg94ktQ1iij15kyqINbiTcIW70jEXwIbfQAb0oI1-GsmeW55Y", "_user_domain_id": "default", "resource_uuid": null, "cell_uuid": null, "service_project_domain_name": null, "read_only": false, "system_scope": null, "service_project_id": null, "domain_name": null, "is_admin_project": true, "service_user_name": null, "user_name": "admin", "user_domain_name": null, "_user_id": "a1e0d70fc2274d39af8434d633c3347e", "project_domain_name": null, "db_connection": null, "project_name": "admin", "global_request_id": "req-609a9b47-1a04-4b23-9646-81e94f86663d", "service_project_name": null, "timestamp": "2023-02-27T06:54:46.917692", "service_project_domain_id": null, "remote_address": "10.249.4.108", "quota_class": null, "_domain_id": null, "user_auth_plugin": null, "service_catalog": [{"endpoints": [{"adminURL": "http://10.50.31.1:8780", "region": "RegionOne", "internalURL": "http://10.50.31.1:8780", "publicURL": "http://10.50.31.1:8780"}], "type": "placement", "name": "placement"}, {"endpoints": [{"adminURL": "http://10.50.31.1:9696", "region": "RegionOne", "internalURL": "http://10.50.31.1:9696", "publicURL": "http://10.50.31.1:9696"}], "type": "network", "name": "neutron"}, {"endpoints": [{"adminURL": "http://10.50.31.1:8776/v3/7e8babd4464e4c6da382a1a29d8da53a", "region": "RegionOne", "internalURL": "http://10.50.31.1:8776/v3/7e8babd4464e4c6da382a1a29d8da53a", "publicURL": "http://10.50.31.1:8776/v3/7e8babd4464e4c6da382a1a29d8da53a"}], "type": "volumev3", "name": "cinderv3"}, {"endpoints": [{"adminURL": "http://10.50.31.1:9292", "region": "RegionOne", "internalURL": "http://10.50.31.1:9292", "publicURL": "http://10.50.31.1:9292"}], "type": "image", "name": "glance"}], "_project_id": "7e8babd4464e4c6da382a1a29d8da53a", "show_deleted": false, "service_roles": [], "service_token": null, "roles": ["admin"], "service_user_domain_id": null, "_read_deleted": "no", "request_id": "req-fb747d3a-977d-4299-baa9-672e37e9d54f", "mq_connection": null, "is_admin": true, "_project_domain_id": "default"}, "_obj_created_at": "2023-02-22T07:06:18.000000", "_obj_deleted": false, "_obj_deleted_at": null}, "hostname": "dx-instance12", "launched_on": "con01.vim1.local", "display_description": null, "key_data": null, "vcpu_model": {"_obj_vendor": null, "_changed_fields": ["vendor", "features", "mode", "model", "arch", "match", "topology"], "VERSION": "1.0", "_obj_features": [], "_obj_model": null, "_obj_mode": "host-model", "_obj_arch": null, "_context": {"service_user_domain_name": null, "service_user_id": null, "auth_token": "gAAAAABj_CrUITIzmjumgkA3iaAICzaDqNNHxfcaA3pF6KQzrOjJjPRsHmy-9UJ52pu2Befpe0GqNfOZpU4GLHJELLTGai8e8YxhHdK-cN2tA2XAU_WhmgO5z8TPrKfeZBPzbo3jRrCwDzoOcE_piPmdu39dg94ktQ1iij15kyqINbiTcIW70jEXwIbfQAb0oI1-GsmeW55Y", "_user_domain_id": "default", "resource_uuid": null, "cell_uuid": null, "service_project_domain_name": null, "read_only": false, "system_scope": null, "service_project_id": null, "domain_name": null, "is_admin_project": true, "service_user_name": null, "user_name": "admin", "user_domain_name": null, "_user_id": "a1e0d70fc2274d39af8434d633c3347e", "project_domain_name": null, "db_connection": null, "project_name": "admin", "global_request_id": "req-609a9b47-1a04-4b23-9646-81e94f86663d", "service_project_name": null, "timestamp": "2023-02-27T06:54:46.917692", "service_project_domain_id": null, "remote_address": "10.249.4.108", "quota_class": null, "_domain_id": null, "user_auth_plugin": null, "service_catalog": [{"endpoints": [{"adminURL": "http://10.50.31.1:8780", "region": "RegionOne", "internalURL": "http://10.50.31.1:8780", "publicURL": "http://10.50.31.1:8780"}], "type": "placement", "name": "placement"}, {"endpoints": [{"adminURL": "http://10.50.31.1:9696", "region": "RegionOne", "internalURL": "http://10.50.31.1:9696", "publicURL": "http://10.50.31.1:9696"}], "type": "network", "name": "neutron"}, {"endpoints": [{"adminURL": "http://10.50.31.1:8776/v3/7e8babd4464e4c6da382a1a29d8da53a", "region": "RegionOne", "internalURL": "http://10.50.31.1:8776/v3/7e8babd4464e4c6da382a1a29d8da53a", "publicURL": "http://10.50.31.1:8776/v3/7e8babd4464e4c6da382a1a29d8da53a"}], "type": "volumev3", "name": "cinderv3"}, {"endpoints": [{"adminURL": "http://10.50.31.1:9292", "region": "RegionOne", "internalURL": "http://10.50.31.1:9292", "publicURL": "http://10.50.31.1:9292"}], "type": "image", "name": "glance"}], "_project_id": "7e8babd4464e4c6da382a1a29d8da53a", "show_deleted": false, "service_roles": [], "service_token": null, "roles": ["admin"], "service_user_domain_id": null, "_read_deleted": "no", "request_id": "req-fb747d3a-977d-4299-baa9-672e37e9d54f", "mq_connection": null, "is_admin": true, "_project_domain_id": "default"}, "_obj_match": "exact", "_obj_topology": {"_changed_fields": ["cores", "threads", "sockets"], "_obj_sockets": 32, "VERSION": "1.0", "_obj_cores": 1, "_obj_threads": 1, "_context": null}}, "power_state": 1, "device_metadata": null, "default_ephemeral_device": null, "migration_context": {"_obj_instance_uuid": "a55c993a-4012-4b41-97d3-69d0666e8197", "_context": {"service_user_domain_name": null, "service_user_id": null, "auth_token": "gAAAAABj_CrUITIzmjumgkA3iaAICzaDqNNHxfcaA3pF6KQzrOjJjPRsHmy-9UJ52pu2Befpe0GqNfOZpU4GLHJELLTGai8e8YxhHdK-cN2tA2XAU_WhmgO5z8TPrKfeZBPzbo3jRrCwDzoOcE_piPmdu39dg94ktQ1iij15kyqINbiTcIW70jEXwIbfQAb0oI1-GsmeW55Y", "_user_domain_id": "default", "resource_uuid": null, "cell_uuid": null, "service_project_domain_name": null, "read_only": false, "system_scope": null, "service_project_id": null, "domain_name": null, "is_admin_project": true, "service_user_name": null, "user_name": "admin", "user_domain_name": null, "_user_id": "a1e0d70fc2274d39af8434d633c3347e", "project_domain_name": null, "db_connection": null, "project_name": "admin", "global_request_id": "req-609a9b47-1a04-4b23-9646-81e94f86663d", "service_project_name": null, "timestamp": "2023-02-27T06:54:46.917692", "service_project_domain_id": null, "remote_address": "10.249.4.108", "quota_class": null, "_domain_id": null, "user_auth_plugin": null, "service_catalog": [{"endpoints": [{"adminURL": "http://10.50.31.1:8780", "region": "RegionOne", "internalURL": "http://10.50.31.1:8780", "publicURL": "http://10.50.31.1:8780"}], "type": "placement", "name": "placement"}, {"endpoints": [{"adminURL": "http://10.50.31.1:9696", "region": "RegionOne", "internalURL": "http://10.50.31.1:9696", "publicURL": "http://10.50.31.1:9696"}], "type": "network", "name": "neutron"}, {"endpoints": [{"adminURL": "http://10.50.31.1:8776/v3/7e8babd4464e4c6da382a1a29d8da53a", "region": "RegionOne", "internalURL": "http://10.50.31.1:8776/v3/7e8babd4464e4c6da382a1a29d8da53a", "publicURL": "http://10.50.31.1:8776/v3/7e8babd4464e4c6da382a1a29d8da53a"}], "type": "volumev3", "name": "cinderv3"}, {"endpoints": [{"adminURL": "http://10.50.31.1:9292", "region": "RegionOne", "internalURL": "http://10.50.31.1:9292", "publicURL": "http://10.50.31.1:9292"}], "type": "image", "name": "glance"}], "_project_id": "7e8babd4464e4c6da382a1a29d8da53a", "show_deleted": false, "service_roles": [], "service_token": null, "roles": ["admin"], "service_user_domain_id": null, "_read_deleted": "no", "request_id": "req-fb747d3a-977d-4299-baa9-672e37e9d54f", "mq_connection": null, "is_admin": true, "_project_domain_id": "default"}, "_obj_new_pci_requests": {"instance_uuid": "a55c993a-4012-4b41-97d3-69d0666e8197", "requests": []}, "_obj_old_resources": null, "_obj_new_resources": null, "_obj_new_numa_topology": null, "_changed_fields": ["new_pci_requests", "old_numa_topology", "old_resources", "instance_uuid", "migration_id", "new_resources", "new_numa_topology", "old_pci_requests", "new_pci_devices", "old_pci_devices"], "VERSION": "1.2", "_obj_new_pci_devices": [], "_obj_old_pci_devices": [], "_obj_old_numa_topology": {"cells": [null], "emulator_threads_policy": null}, "_obj_migration_id": 699, "_obj_old_pci_requests": {"instance_uuid": "a55c993a-4012-4b41-97d3-69d0666e8197", "requests": []}}, "hidden": false, "project_id": "7e8babd4464e4c6da382a1a29d8da53a", "launched_at": "2023-02-27T06:47:41.000000", "resources": null, "config_drive": "True", "node": "con01.vim1.local", "pci_devices": [], "access_ip_v6": null, "access_ip_v4": null, "deleted": false, "key_name": null, "host": "con01.vim1.local", "ephemeral_key_uuid": null, "progress": 0, "services": [{"binary": "nova-compute", "uuid": "936011a4-4581-4b53-8a76-415a3bd4da7b", "deleted": false, "created_at": "2021-10-09T04:13:42.000000", "updated_at": "2023-02-27T06:55:57.000000", "report_count": 4363096, "topic": "compute", "host": "con01.vim1.local", "version": 40, "disabled": false, "forced_down": false, "last_seen_up": "2023-02-27T06:55:57.000000", "deleted_at": null, "disabled_reason": null, "id": 6}], "display_name": "dx_instance12", "system_metadata": {"image_os_distro": "centos", "owner_user_name": "admin", "image_os_admin_user": "root", "image_image_type": "image", "image_os_version": "11", "boot_roles": "admin", "clean_attempts": "1", "image_disk_format": "raw", "image_hw_qemu_guest_agent": "yes", "image_container_format": "bare", "image_min_ram": "0", "old_vm_state": "active", "image_min_disk": "0", "image_usage_type": "common", "owner_project_name": "admin"}, "task_state": null, "shutdown_terminate": false, "os_type": null, "cell_name": null, "root_gb": 0, "kernel_id": "", "name": "instance-00002db0", "instance_type_id": 801, "locked_by": null, "launch_index": 0, "locked": false, "memory_mb": 4096, "vcpus": 1, "image_ref": "a9cd271f-93b8-4e08-9a19-e328c46b1ac4", "root_device_name": "/dev/vda", "auto_disk_config": false, "new_flavor": {"memory_mb": 131072, "root_gb": 0, "deleted_at": null, "description": null, "deleted": false, "created_at": "2022-11-01T07:05:54.000000", "ephemeral_gb": 0, "updated_at": null, "disabled": false, "vcpus": 32, "extra_specs": {}, "swap": 0, "rxtx_factor": 1.0, "is_public": true, "flavorid": "03a4fe3d-2b87-466d-88d6-430ec5ea0637", "vcpu_weight": 0, "id": 951, "name": "CPU\u4e91\u670d\u52a1\u5668_VCS.CGI.32A_1667286354021"}, "architecture": null, "metadata": {}, "ramdisk_id": "", "created_at": "2023-02-22T07:06:17.000000"}}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("finish_resize", "nova", "nova", msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != "ironic" || publisher.ErrorMsg == nil ||
		publisher.ResourceId != "a55c993a-4012-4b41-97d3-69d0666e8197" ||
		publisher.Host != "con01.vim1.local" || publisher.ResourceState != "error" ||
		publisher.FlavorId != "0348236c-417d-4409-87db-36550fdeebe8" {
		t.Fatal("Failed to handleNova")
	}
}

func TestHandleNova_compute_task_build_instances(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"request_spec": {"instance_properties": {"root_gb": 0, "user_id": "83bac5ab80b14f74b06f157daa6fd64b", "uuid": "bfeabbda-cc1d-45f5-b78e-5b6cd235f709", "availability_zone": "nova", "ephemeral_gb": 0, "numa_topology": {"cells": [{"cpuset_reserved": null, "pagesize": -3, "cpuset": [0, 1, 2, 3], "cpu_policy": null, "memory": 8192, "cpu_pinning_raw": null, "id": 0, "cpu_thread_policy": null}]}, "memory_mb": 8192, "vcpus": 4, "project_id": "a01edbe369764a6ba25798bb477f245b", "pci_requests": {"requests": []}}, "instance_type": {"memory_mb": 8192, "root_gb": 0, "deleted_at": null, "description": null, "deleted": false, "created_at": "2022-09-28T09:56:04.000000", "ephemeral_gb": 0, "updated_at": null, "disabled": false, "vcpus": 4, "extra_specs": {":category": "general_purpose", "hw:mem_page_size": "any", "hw:live_resize": "True", ":architecture": "x86_architecture", "hw:numa_nodes": "1"}, "swap": 0, "rxtx_factor": 1.0, "is_public": true, "flavorid": "0fa48b8a-ccd1-41ce-828a-56d867bb90f9", "vcpu_weight": 0, "id": 3, "name": "SEPC"}, "image": {"status": "active", "properties": {"os_distro": "centos", "hw_qemu_guest_agent": true, "os_type": "linux", "os_admin_user": "root"}, "name": "centos7.9-220929", "container_format": "bare", "created_at": "2022-09-29T02:39:40.000000", "disk_format": "raw", "updated_at": "2022-09-29T02:41:32.000000", "id": "032e83c9-eddb-40d8-8c23-e4e8015c79ee", "owner": "a01edbe369764a6ba25798bb477f245b", "checksum": "20484673cba6cbd822ba71932af9fb43", "min_disk": 0, "min_ram": 0, "size": 3221225472}, "num_instances": 1}, "reason": "{'message': u'Exceeded maximum number of retries. Exhausted all hosts available for retrying build failures for instance bfeabbda-cc1d-45f5-b78e-5b6cd235f709.', 'class': 'MaxRetriesExceeded', 'kwargs': {'reason': u'Exhausted all hosts available for retrying build failures for instance bfeabbda-cc1d-45f5-b78e-5b6cd235f709.', 'code': 500}}", "instance_id": "bfeabbda-cc1d-45f5-b78e-5b6cd235f709", "state": "error", "instance_properties": {"root_gb": 0, "user_id": "83bac5ab80b14f74b06f157daa6fd64b", "uuid": "bfeabbda-cc1d-45f5-b78e-5b6cd235f709", "availability_zone": "nova", "ephemeral_gb": 0, "numa_topology": {"cells": [{"cpuset_reserved": null, "pagesize": -3, "cpuset": [0, 1, 2, 3], "cpu_policy": null, "memory": 8192, "cpu_pinning_raw": null, "id": 0, "cpu_thread_policy": null}]}, "memory_mb": 8192, "vcpus": 4, "project_id": "a01edbe369764a6ba25798bb477f245b", "pci_requests": {"requests": []}}, "method": "build_instances"}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("compute_task.build_instances", "nova", "nova", msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != "ironic" || publisher.ErrorMsg == nil ||
		publisher.ResourceId != "bfeabbda-cc1d-45f5-b78e-5b6cd235f709" ||
		publisher.Host != "" || publisher.ResourceState != "error" ||
		publisher.FlavorId != "0fa48b8a-ccd1-41ce-828a-56d867bb90f9"{
		t.Fatal("Failed to handleNova")
	}
}

func TestHandleNova_compute_task_migrate_server(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"request_spec": {"instance_properties": {"root_gb": 0, "user_id": "a1e0d70fc2274d39af8434d633c3347e", "uuid": "416a0a50-f20e-4308-b326-edfc67faf203", "availability_zone": "nova", "ephemeral_gb": 0, "numa_topology": null, "memory_mb": 131072, "vcpus": 32, "project_id": "7e8babd4464e4c6da382a1a29d8da53a", "pci_requests": {"requests": []}}, "instance_type": {"memory_mb": 131072, "root_gb": 0, "deleted_at": null, "description": null, "deleted": false, "created_at": "2022-11-01T07:05:54.000000", "ephemeral_gb": 0, "updated_at": null, "disabled": false, "vcpus": 32, "extra_specs": {}, "swap": 0, "rxtx_factor": 1.0, "is_public": true, "flavorid": "03a4fe3d-2b87-466d-88d6-430ec5ea0637", "vcpu_weight": 0, "id": 951, "name": "CPU\u4e91\u670d\u52a1\u5668_VCS.CGI.32A_1667286354021"}, "image": {"status": "active", "properties": {"os_distro": "centos", "hw_qemu_guest_agent": true, "os_admin_user": "root"}, "name": "centos7", "container_format": "bare", "created_at": "2021-10-09T13:18:34.000000", "disk_format": "raw", "updated_at": "2021-11-01T06:57:47.000000", "id": "35a7999f-e43d-40b4-9ea6-e0cd0d67ea55", "owner": "7e8babd4464e4c6da382a1a29d8da53a", "checksum": "514ce7dea76fce7ee5e3cf2e68ef91fb", "min_disk": 10, "min_ram": 1024, "size": 3221225472}, "num_instances": 1}, "reason": "{'message': u'No valid host was found. ', 'class': 'NoValidHost_Remote', 'kwargs': {u'reason': u'', u'code': 500}}", "instance_id": "416a0a50-f20e-4308-b326-edfc67faf203", "state": "stopped", "instance_properties": {"root_gb": 0, "user_id": "a1e0d70fc2274d39af8434d633c3347e", "uuid": "416a0a50-f20e-4308-b326-edfc67faf203", "availability_zone": "nova", "ephemeral_gb": 0, "numa_topology": null, "memory_mb": 131072, "vcpus": 32, "project_id": "7e8babd4464e4c6da382a1a29d8da53a", "pci_requests": {"requests": []}}, "method": "migrate_server"}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("compute_task.migrate_server", consts.NOVA, consts.NOVA, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.INSTANCE || publisher.ErrorMsg == nil ||
		publisher.ResourceId != "416a0a50-f20e-4308-b326-edfc67faf203" ||
		publisher.Host != "" || publisher.ResourceState != "stopped" ||
		publisher.FlavorId != "03a4fe3d-2b87-466d-88d6-430ec5ea0637" {
		t.Fatal("Failed to handleNova")
	}
}

func TestHandleNova_compute_instance_resize_error(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"state_description": "resize_prep", "code": 500, "availability_zone": "nova", "terminated_at": "", "ephemeral_gb": 0, "instance_type_id": 4, "message": "No valid host was found. No valid host found for cold migrate", "deleted_at": "", "reservation_id": "r-ij31aj0v", "instance_id": "ed90e52c-1a55-46c2-95e0-48d7bedf66c6", "display_name": "hanjuntao_dev4444-1730120357155999745", "hostname": "hanjuntao-dev4444-1730120357155999745", "state": "stopped", "progress": "", "launched_at": "2023-11-30T07:08:05.000000", "metadata": {}, "node": "allinone-02", "ramdisk_id": "", "access_ip_v6": null, "disk_gb": 0, "access_ip_v4": null, "kernel_id": "", "host": "allinone-02", "user_id": "d7fb04a0db0042e6b663c1627d90a383", "image_ref_url": "http://10.251.28.249:9292/images/1e08106a-a25c-4565-b1ad-87484ff95608", "cell_name": "", "exception": "{'message': u'No valid host was found. No valid host found for cold migrate', 'class': 'NoValidHost_Remote', 'kwargs': {u'reason': u'No valid host found for cold migrate', u'code': 500}}", "root_gb": 0, "tenant_id": "0dfc4c30add9460d9dbb1907a4af8032", "created_at": "2023-11-30 07:07:52+00:00", "memory_mb": 2048, "instance_type": "SECS_\u901a\u7528\u578bsecs.g1.small.2", "vcpus": 1, "image_meta": {"os_distro": "centos", "image_type": "image", "container_format": "bare", "hw_qemu_guest_agent": "yes", "disk_format": "raw", "os_admin_user": "root", "base_image_ref": "1e08106a-a25c-4565-b1ad-87484ff95608", "os_version": "7.9", "min_disk": "20", "usage_type": "common", "min_ram": "0"}, "architecture": null, "os_type": null, "instance_flavor_id": "d3ad430e-c610-4c38-8434-cd3efbef404f"}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("compute.instance.resize.error", consts.NOVA, consts.NOVA, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.INSTANCE || publisher.ErrorMsg != "{'message': u'No valid host was found. No valid host found for cold migrate', 'class': 'NoValidHost_Remote', 'kwargs': {u'reason': u'No valid host found for cold migrate', u'code': 500}}" ||
		publisher.ResourceId != "ed90e52c-1a55-46c2-95e0-48d7bedf66c6" ||
		publisher.Host != "allinone-02" || publisher.ResourceState != "stopped" ||
		publisher.FlavorId != "d3ad430e-c610-4c38-8434-cd3efbef404f" {
		t.Fatal("Failed to handleNova")
	}
}

// BUG http://wiki.voneyun.com/pages/viewpage.action?pageId=91172785
func TestHandleNova_compute_instance_resize_error_retry(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"state_description": "resize_prep", "code": 500, "availability_zone": "nova", "terminated_at": "", "ephemeral_gb": 0, "instance_type_id": 4, "message": "{'class': 'InstanceFaultRollback', 'message': u'Instance rollback performed due to: Unable to migrate instance (248480ee-2f96-4785-a5de-6499874ee834) to current host (SCDA0108.uat.local).', 'inner_exception': {'message': u'Unable to migrate instance (248480ee-2f96-4785-a5de-6499874ee834) to current host (SCDA0108.uat.local).', 'class': 'UnableToMigrateToSelf', 'kwargs': {'instance_id': u'248480ee-2f96-4785-a5de-6499874ee834', 'host': 'SCDA0108.uat.local', 'code': 400}}, 'kwargs': {'code': 500}}", "deleted_at": "", "reservation_id": "r-ij31aj0v", "instance_id": "ed90e52c-1a55-46c2-95e0-48d7bedf66c6", "display_name": "hanjuntao_dev4444-1730120357155999745", "hostname": "hanjuntao-dev4444-1730120357155999745", "state": "stopped", "progress": "", "launched_at": "2023-11-30T07:08:05.000000", "metadata": {}, "node": "allinone-02", "ramdisk_id": "", "access_ip_v6": null, "disk_gb": 0, "access_ip_v4": null, "kernel_id": "", "host": "allinone-02", "user_id": "d7fb04a0db0042e6b663c1627d90a383", "image_ref_url": "http://10.251.28.249:9292/images/1e08106a-a25c-4565-b1ad-87484ff95608", "cell_name": "", "exception": "{'class': 'InstanceFaultRollback', 'message': u'Instance rollback performed due to: Unable to migrate instance (248480ee-2f96-4785-a5de-6499874ee834) to current host (SCDA0108.uat.local).', 'inner_exception': {'message': u'Unable to migrate instance (248480ee-2f96-4785-a5de-6499874ee834) to current host (SCDA0108.uat.local).', 'class': 'UnableToMigrateToSelf', 'kwargs': {'instance_id': u'248480ee-2f96-4785-a5de-6499874ee834', 'host': 'SCDA0108.uat.local', 'code': 400}}, 'kwargs': {'code': 500}}", "root_gb": 0, "tenant_id": "0dfc4c30add9460d9dbb1907a4af8032", "created_at": "2023-11-30 07:07:52+00:00", "memory_mb": 2048, "instance_type": "SECS_\u901a\u7528\u578bsecs.g1.small.2", "vcpus": 1, "image_meta": {"os_distro": "centos", "image_type": "image", "container_format": "bare", "hw_qemu_guest_agent": "yes", "disk_format": "raw", "os_admin_user": "root", "base_image_ref": "1e08106a-a25c-4565-b1ad-87484ff95608", "os_version": "7.9", "min_disk": "20", "usage_type": "common", "min_ram": "0"}, "architecture": null, "os_type": null, "instance_flavor_id": "d3ad430e-c610-4c38-8434-cd3efbef404f"}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("compute.instance.resize.error", consts.NOVA, consts.NOVA, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher != nil {
		t.Fatal("Failed to handleNova")
	}
}

func TestHandleNova_compute_instance_resize_confirm_end(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"state_description": "", "availability_zone": "nova", "terminated_at": "", "ephemeral_gb": 0, "instance_type_id": 1069, "deleted_at": "", "fixed_ips": [{"version": 4, "vif_mac": "fa:16:3e:dd:c8:48", "floating_ips": [], "label": "dx_vxlan1", "meta": {}, "address": "12.12.12.15", "type": "fixed"}], "instance_id": "416a0a50-f20e-4308-b326-edfc67faf203", "display_name": "test", "reservation_id": "r-x481qy1i", "hostname": "test", "state": "stopped", "progress": "", "launched_at": "2023-04-28T02:27:26.000000", "metadata": {}, "node": "con01.vim1.local", "ramdisk_id": "", "access_ip_v6": null, "disk_gb": 0, "access_ip_v4": null, "kernel_id": "", "host": "con01.vim1.local", "user_id": "a1e0d70fc2274d39af8434d633c3347e", "image_ref_url": "http://10.50.31.1:9292/images/35a7999f-e43d-40b4-9ea6-e0cd0d67ea55", "cell_name": "", "root_gb": 0, "tenant_id": "7e8babd4464e4c6da382a1a29d8da53a", "created_at": "2023-04-28 02:24:36+00:00", "memory_mb": 8192, "instance_type": "004008", "vcpus": 4, "image_meta": {"os_distro": "centos", "image_type": "image", "container_format": "bare", "hw_qemu_guest_agent": "yes", "disk_format": "raw", "os_admin_user": "root", "os_version": "7.6", "min_ram": "1024", "min_disk": "10", "usage_type": "common", "base_image_ref": "35a7999f-e43d-40b4-9ea6-e0cd0d67ea55"}, "architecture": null, "os_type": null, "instance_flavor_id": "004008"}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("compute.instance.resize.confirm.end", consts.NOVA, consts.NOVA, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.INSTANCE || publisher.ErrorMsg != nil ||
		publisher.ResourceId != "416a0a50-f20e-4308-b326-edfc67faf203" ||
		publisher.Host != "con01.vim1.local" || publisher.ResourceState != "stopped" ||
		publisher.FlavorId != "004008" {
		t.Fatal("Failed to handleNova")
	}
}

func TestHandleNova_live_migration_rollback_dest_end(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"state_description": "", "availability_zone": "amd", "terminated_at": "", "ephemeral_gb": 0, "instance_type_id": 790, "deleted_at": "", "fixed_ips": [{"version": 4, "vif_mac": "fa:16:3e:a7:6b:34", "floating_ips": [], "label": "sincere_vpc777", "meta": {}, "address": "10.5.0.131", "type": "fixed"}], "instance_id": "f5ac7367-89b6-4202-a7a8-4ab9b3f4e503", "display_name": "shaohaoyu-1674363344822435842", "reservation_id": "r-o2cmpjow", "hostname": "shaohaoyu-1674363344822435842", "state": "active", "progress": "", "launched_at": "2023-06-29T10:25:17.000000", "metadata": {}, "node": "SCDA0072.uat.local", "ramdisk_id": "", "access_ip_v6": null, "disk_gb": 0, "access_ip_v4": null, "kernel_id": "", "host": "SCDA0072.uat.local", "user_id": "f8d81f2862e44a60aac5bc407ce6889d", "image_ref_url": "http://10.250.4.200:9292/images/e8f2a1df-ee0b-45fe-8bfc-3b9f07eab506", "cell_name": "", "root_gb": 0, "tenant_id": "fa8d250efe314a0fb8daa197bd82e9be", "created_at": "2023-06-29 10:25:06+00:00", "memory_mb": 1024, "instance_type": "SECS_\u4e91\u670d\u52a1\u5668-\u901a\u7528\u578bC2secs.c2.small.1", "vcpus": 1, "image_meta": {"hw_qemu_guest_agent": "yes", "os_distro": "centos", "image_type": "image", "container_format": "bare", "min_ram": "0", "owner_specified.openstack.sha256": "", "disk_format": "raw", "os_admin_user": "root", "usage_type": "common", "os_version": "8.0", "owner_specified.openstack.object": "images/Centos8.0", "owner_specified.openstack.md5": "", "min_disk": "20", "os_type": "linux", "base_image_ref": "e8f2a1df-ee0b-45fe-8bfc-3b9f07eab506"}, "architecture": null, "os_type": "linux", "instance_flavor_id": "fa476411-a99e-4ac0-80df-82c8033984cc"}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("compute.instance.live_migration.rollback.dest.end", consts.NOVA, consts.NOVA, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.INSTANCE || publisher.ErrorMsg != nil ||
		publisher.ResourceId != "f5ac7367-89b6-4202-a7a8-4ab9b3f4e503" ||
		publisher.Host != "SCDA0072.uat.local" || publisher.Node != "SCDA0072.uat.local" ||
		!reflect.DeepEqual(publisher.FixedIps, []string{"10.5.0.131"}) ||
		publisher.ResourceState != "active" || publisher.FlavorId != "fa476411-a99e-4ac0-80df-82c8033984cc" {
		t.Fatal("Failed to handleNova")
	}
}

// instance live resize end
func TestHandleNeutron_instance_live_resize_end(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"nova_object.version": "1.8", "nova_object.name": "InstanceActionPayload", "nova_object.namespace": "nova", "nova_object.data": {"availability_zone": "nova", "terminated_at": null, "ip_addresses": [{"nova_object.version": "1.0", "nova_object.name": "IpPayload", "nova_object.namespace": "nova", "nova_object.data": {"port_uuid": "ef6020d7-0668-41d2-98e0-0bc2467d9270", "device_name": "tapef6020d7-06", "mac": "fa:16:3e:0c:d4:62", "version": 4, "meta": {}, "address": "192.101.101.100", "label": "sdn_test33696160203935745_network"}}], "ramdisk_id": "", "updated_at": "2023-11-30T02:35:34Z", "image_uuid": "1e08106a-a25c-4565-b1ad-87484ff95608", "flavor": {"nova_object.version": "1.4", "nova_object.name": "FlavorPayload", "nova_object.namespace": "nova", "nova_object.data": {"memory_mb": 8192, "root_gb": 0, "description": null, "ephemeral_gb": 0, "disabled": false, "vcpus": 2, "extra_specs": {":category": "general_purpose", "hw:mem_page_size": "any", ":architecture": "x86_architecture", "hw:live_resize": "True", "hw:numa_nodes": "1"}, "swap": 0, "rxtx_factor": 1.0, "is_public": true, "flavorid": "a433ee13-c196-4935-8119-bf8981474ae6", "vcpu_weight": 0, "projects": null, "name": "c2s8"}}, "deleted_at": null, "reservation_id": "r-9xzqkqbk", "display_name": "sdn_test33696848325644289_server", "uuid": "73cb238b-c37a-467f-a630-53cfda2c9ca4", "display_description": "sdn_test33696848325644289_server", "action_initiator_user": "8ea8f37c546c46a68ccfda71de73d4a9", "locked_reason": null, "state": "active", "power_state": "running", "host_name": "sdn-test33696848325644289-server", "progress": 0, "launched_at": "2023-11-29T08:18:48Z", "metadata": {}, "node": "allinone-02", "action_initiator_project": "e97b968db6324b9c9f68d9dc70af3462", "kernel_id": "", "key_name": null, "host": "allinone-02", "user_id": "8ea8f37c546c46a68ccfda71de73d4a9", "task_state": null, "locked": false, "tenant_id": "e97b968db6324b9c9f68d9dc70af3462", "created_at": "2023-11-29T08:18:36Z", "block_devices": [{"nova_object.version": "1.0", "nova_object.name": "BlockDevicePayload", "nova_object.namespace": "nova", "nova_object.data": {"device_name": "/dev/vda", "boot_index": 0, "tag": null, "delete_on_termination": true, "volume_id": "70df3684-0f08-4acf-9663-a6ef9127de52"}}], "architecture": null, "request_id": "req-dba3a143-32c7-46ff-bcb9-8763e0cd7248", "auto_disk_config": "MANUAL", "os_type": null, "fault": null}}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("instance.live_resize.end", consts.NOVA, consts.NOVA, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.INSTANCE || publisher.ErrorMsg != nil ||
		publisher.ResourceId != "73cb238b-c37a-467f-a630-53cfda2c9ca4" ||
		publisher.ResourceState != "active" || publisher.Host != "allinone-02" ||
		publisher.Node != "allinone-02" || publisher.FlavorId != "a433ee13-c196-4935-8119-bf8981474ae6"{
		t.Fatal("Failed to handleNeutron")
	}
}

// instance live resize error
func TestHandleNeutron_instance_live_resize_error(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"request_spec": {"instance_properties": {"root_gb": 0, "user_id": "d7fb04a0db0042e6b663c1627d90a383", "uuid": "f654c5a4-1ca1-4fdb-824f-dbebe696e3ab", "availability_zone": "nova", "ephemeral_gb": 0, "numa_topology": {"cells": [{"cpuset_reserved": null, "pagesize": 4, "cpuset": [0], "cpu_policy": null, "memory": 1024, "cpu_pinning_raw": null, "id": 0, "cpu_thread_policy": null}], "emulator_threads_policy": null}, "memory_mb": 16384, "vcpus": 12, "project_id": "0dfc4c30add9460d9dbb1907a4af8032", "pci_requests": {"instance_uuid": "f654c5a4-1ca1-4fdb-824f-dbebe696e3ab", "requests": []}}, "instance_type": {"memory_mb": 16384, "root_gb": 0, "deleted_at": null, "description": null, "deleted": false, "created_at": "2023-11-29T03:30:44.000000", "ephemeral_gb": 0, "updated_at": null, "disabled": false, "vcpus": 12, "extra_specs": {"hw:mem_page_size": "any", "hw:live_resize": "True", ":category": "general_purpose", "ecs": "amd", "hw:numa_nodes": "1", ":architecture": "x86_architecture"}, "swap": 0, "rxtx_factor": 1.0, "is_public": true, "flavorid": "e31fe9e7-b7e3-46a5-8f68-0a3aeaf39dc4", "vcpu_weight": 0, "id": 6, "name": "SECS_\u901a\u7528\u578bsecs.g1.small.3"}, "image": {"min_disk": 20, "container_format": "bare", "min_ram": 0, "disk_format": "raw", "properties": {"os_distro": "centos", "hw_qemu_guest_agent": true, "os_admin_user": "root"}}, "num_instances": 1}, "reason": "{'message': u'No valid host was found. There are not enough hosts available.', 'class': 'NoValidHost_Remote', 'kwargs': {u'reason': u'There are not enough hosts available.', u'code': 500}}", "instance_id": "f654c5a4-1ca1-4fdb-824f-dbebe696e3ab", "state": "error", "instance_properties": {"root_gb": 0, "user_id": "d7fb04a0db0042e6b663c1627d90a383", "uuid": "f654c5a4-1ca1-4fdb-824f-dbebe696e3ab", "availability_zone": "nova", "ephemeral_gb": 0, "numa_topology": {"cells": [{"cpuset_reserved": null, "pagesize": 4, "cpuset": [0], "cpu_policy": null, "memory": 1024, "cpu_pinning_raw": null, "id": 0, "cpu_thread_policy": null}], "emulator_threads_policy": null}, "memory_mb": 16384, "vcpus": 12, "project_id": "0dfc4c30add9460d9dbb1907a4af8032", "pci_requests": {"instance_uuid": "f654c5a4-1ca1-4fdb-824f-dbebe696e3ab", "requests": []}}, "method": "live_resize"}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("compute_task.live_resize", consts.NOVA, consts.NOVA, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.INSTANCE ||
		publisher.ResourceId != "f654c5a4-1ca1-4fdb-824f-dbebe696e3ab" ||
		publisher.ResourceState != "error" || publisher.Host != "" ||
		publisher.FlavorId != "e31fe9e7-b7e3-46a5-8f68-0a3aeaf39dc4" ||
		publisher.ErrorMsg != "{'message': u'No valid host was found. There are not enough hosts available.', 'class': 'NoValidHost_Remote', 'kwargs': {u'reason': u'There are not enough hosts available.', u'code': 500}}" {
		t.Fatal("Failed to handleNeutron")
	}
}

func TestHandleCompute_compute_service(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"nova_object.version": "1.1", "nova_object.name": "ServiceStatusPayload", "nova_object.namespace": "nova", "nova_object.data": {"binary": "nova-compute", "uuid": "936011a4-4581-4b53-8a76-415a3bd4da7b", "availability_zone": null, "report_count": 4711304, "topic": "compute", "host": "con01.vim1.local", "version": 40, "disabled": false, "forced_down": false, "last_seen_up": "2023-04-11T09:42:40Z", "disabled_reason": "null"}}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("service.update", consts.ComputeService, consts.COMPUTE, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.ComputeService || publisher.ErrorMsg != "null" ||
		publisher.ResourceId != "936011a4-4581-4b53-8a76-415a3bd4da7b" ||
		publisher.Host != "con01.vim1.local" || publisher.ResourceState != "" || publisher.ProvisionState != "enabled"{
		t.Fatal("Failed to handleCompute")
	}
}

func TestHandleCinder_attach_volume(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"exception": "{'message': u'Instance a73f2dec-8dca-4ae0-a113-c21778f61a4a could not be found.', 'class': 'InstanceNotFound', 'kwargs': {'instance_id': u'a73f2dec-8dca-4ae0-a113-c21778f61a4a', 'code': 404}}", "args": {"bdm": {"guest_format": null, "attachment_id": "ff988445-d355-4290-a596-c46fd4b28e4d", "updated_at": "2023-04-14T07:12:13.000000", "tag": null, "device_type": null, "snapshot_id": null, "deleted_at": null, "id": 17470, "uuid": "b499897b-340d-40d9-a5de-8f982975111b", "no_device": false, "volume_size": 10, "connection_info": "{}", "destination_type": "volume", "delete_on_termination": false, "boot_index": null, "deleted": false, "image_id": null, "volume_id": "74fbef54-f375-4c75-b058-8d41f850c3ec", "instance_uuid": "a73f2dec-8dca-4ae0-a113-c21778f61a4a", "source_type": "volume", "created_at": "2023-04-14T07:12:12.000000", "volume_type": null, "device_name": "/dev/vdc", "disk_bus": null}, "instance": {"vm_state": "stopped", "availability_zone": "nova", "terminated_at": null, "ephemeral_gb": 0, "instance_type_id": 862, "updated_at": "2023-04-12T07:37:50.000000", "cleaned": true, "vm_mode": null, "deleted_at": null, "reservation_id": "r-mi9ru952", "id": 11969, "security_groups": [], "disable_terminate": false, "user_id": "d113fe136b044262b8e711f6fbecba15", "uuid": "a73f2dec-8dca-4ae0-a113-c21778f61a4a", "default_swap_device": null, "info_cache": {"_obj_instance_uuid": "a73f2dec-8dca-4ae0-a113-c21778f61a4a", "_changed_fields": [], "_obj_updated_at": "2023-04-07T08:51:27.000000", "VERSION": "1.5", "_obj_network_info": [{"profile": {}, "ovs_interfaceid": "2e44ebfb-2e8b-4b68-bce2-67bdf270c75c", "preserve_on_delete": false, "network": {"bridge": "br-int", "label": "dx_net11", "meta": {"injected": false, "tunneled": true, "tenant_id": "e4a97e45239340b090fcebde8edcbfe5", "physical_network": null, "mtu": 1450}, "id": "4b07154f-1cd3-491f-a916-e70fd2fb097b", "subnets": [{"ips": [{"meta": {}, "type": "fixed", "version": 4, "address": "44.44.44.178", "floating_ips": [{"meta": {}, "type": "floating", "version": 4, "address": "10.50.35.242"}]}], "version": 4, "meta": {"dhcp_server": "44.44.44.2"}, "dns": [], "routes": [], "cidr": "44.44.44.0/24", "gateway": {"meta": {}, "type": "gateway", "version": 4, "address": "44.44.44.1"}}]}, "devname": "tap2e44ebfb-2e", "qbh_params": null, "vnic_type": "normal", "meta": {}, "details": {"ovs_hybrid_plug": false, "bridge_name": "br-int", "datapath_type": "system", "port_filter": true, "connectivity": "l2"}, "address": "fa:16:3e:5a:07:1c", "active": false, "type": "ovs", "id": "2e44ebfb-2e8b-4b68-bce2-67bdf270c75c", "qbg_params": null}], "_context": {"service_user_domain_name": null, "service_user_id": null, "auth_token": "gAAAAABkOPws5F54FfNMp7h5pSSdzJPHFLMQGjVrknmxkYSh9M6W84LrKEL6R45PWbi9yLJXtU87m8BMBu4IQzOC7U0ofy6IXI3eAbndN72xKkyZIhhdtQgAaF8QRrEhF9nci48DEfV1a7yy81K2F2bhD3bHqVu1cALueEHPhrvDuTkPm-MsJ_sVgMkpcXI0LJxlTMry0Ai1", "_user_domain_id": "default", "resource_uuid": null, "cell_uuid": null, "service_project_domain_name": null, "read_only": false, "system_scope": null, "service_project_id": null, "domain_name": null, "is_admin_project": true, "service_user_name": null, "user_name": "dx", "user_domain_name": null, "_user_id": "d113fe136b044262b8e711f6fbecba15", "project_domain_name": null, "db_connection": null, "project_name": "dx_project", "global_request_id": "req-c3aa48ff-2a93-4ba0-89aa-85f29bc02883", "service_project_name": null, "timestamp": "2023-04-14T07:12:11.786550", "service_project_domain_id": null, "remote_address": "10.249.4.108", "quota_class": null, "_domain_id": null, "user_auth_plugin": null, "service_catalog": [{"endpoints": [{"adminURL": "http://10.50.31.1:8780", "region": "RegionOne", "internalURL": "http://10.50.31.1:8780", "publicURL": "http://10.50.31.1:8780"}], "type": "placement", "name": "placement"}, {"endpoints": [{"adminURL": "http://10.50.31.1:9696", "region": "RegionOne", "internalURL": "http://10.50.31.1:9696", "publicURL": "http://10.50.31.1:9696"}], "type": "network", "name": "neutron"}, {"endpoints": [{"adminURL": "http://10.50.31.1:8776/v3/e4a97e45239340b090fcebde8edcbfe5", "region": "RegionOne", "internalURL": "http://10.50.31.1:8776/v3/e4a97e45239340b090fcebde8edcbfe5", "publicURL": "http://10.50.31.1:8776/v3/e4a97e45239340b090fcebde8edcbfe5"}], "type": "volumev3", "name": "cinderv3"}, {"endpoints": [{"adminURL": "http://10.50.31.1:9292", "region": "RegionOne", "internalURL": "http://10.50.31.1:9292", "publicURL": "http://10.50.31.1:9292"}], "type": "image", "name": "glance"}], "_project_id": "e4a97e45239340b090fcebde8edcbfe5", "show_deleted": false, "service_roles": [], "service_token": null, "roles": ["system_admin", "courier_system_admin", "ironic_system_admin", "keystone_system_admin", "nova_system_admin", "octavia_system_admin", "panko_system_admin", "neutron_system_admin", "glance_system_admin", "cinder_system_admin", "heat_system_admin", "placement_system_admin"], "service_user_domain_id": null, "_read_deleted": "no", "request_id": "req-c355ded0-22b9-4817-9de0-4b824bc1ecf6", "mq_connection": null, "is_admin": true, "_project_domain_id": "default"}, "_obj_created_at": "2023-03-24T07:09:18.000000", "_obj_deleted": false, "_obj_deleted_at": null}, "hostname": "dx-instance12", "launched_on": "con01.vim1.local", "display_description": null, "key_data": null, "deleted": false, "power_state": 0, "default_ephemeral_device": null, "progress": 0, "hidden": false, "project_id": "e4a97e45239340b090fcebde8edcbfe5", "launched_at": "2023-03-30T09:50:35.000000", "config_drive": "True", "node": "con01.vim1.local", "ramdisk_id": "", "access_ip_v6": null, "access_ip_v4": null, "kernel_id": "", "key_name": null, "user_data": "Q29udGVudC1UeXBlOiBtdWx0aXBhcnQvbWl4ZWQ7IGJvdW5kYXJ5PSI9PT09PT09PT09PT09PT0yMzA5OTg0MDU5NzQzNzYyNDc1PT0iIApNSU1FLVZlcnNpb246IDEuMAoKLS09PT09PT09PT09PT09PT0yMzA5OTg0MDU5NzQzNzYyNDc1PT0KQ29udGVudC1UeXBlOiB0ZXh0L2Nsb3VkLWNvbmZpZzsgY2hhcnNldD0idXMtYXNjaWkiIApNSU1FLVZlcnNpb246IDEuMApDb250ZW50LVRyYW5zZmVyLUVuY29kaW5nOiA3Yml0CkNvbnRlbnQtRGlzcG9zaXRpb246IGF0dGFjaG1lbnQ7IGZpbGVuYW1lPSJzc2gtcHdhdXRoLXNjcmlwdC50eHQiIAoKI2Nsb3VkLWNvbmZpZwpkaXNhYmxlX3Jvb3Q6IGZhbHNlCnNzaF9wd2F1dGg6IHRydWUKcGFzc3dvcmQ6IERvbmd4aWFuZzEyMDcKCi0tPT09PT09PT09PT09PT09MjMwOTk4NDA1OTc0Mzc2MjQ3NT09CkNvbnRlbnQtVHlwZTogdGV4dC94LXNoZWxsc2NyaXB0OyBjaGFyc2V0PSJ1cy1hc2NpaSIgCk1JTUUtVmVyc2lvbjogMS4wCkNvbnRlbnQtVHJhbnNmZXItRW5jb2Rpbmc6IDdiaXQKQ29udGVudC1EaXNwb3NpdGlvbjogYXR0YWNobWVudDsgZmlsZW5hbWU9InBhc3N3ZC1zY3JpcHQudHh0IiAKCiMhL2Jpbi9zaAplY2hvICdyb290OkRvbmd4aWFuZzEyMDcnIHwgY2hwYXNzd2QKCi0tPT09PT09PT09PT09PT09MjMwOTk4NDA1OTc0Mzc2MjQ3NT09CkNvbnRlbnQtVHlwZTogdGV4dC94LXNoZWxsc2NyaXB0OyBjaGFyc2V0PSJ1cy1hc2NpaSIgCk1JTUUtVmVyc2lvbjogMS4wCkNvbnRlbnQtVHJhbnNmZXItRW5jb2Rpbmc6IDdiaXQKQ29udGVudC1EaXNwb3NpdGlvbjogYXR0YWNobWVudDsgZmlsZW5hbWU9ImVuYWJsZS1mcy1jb2xsZWN0b3IudHh0IiAKCiMhL2Jpbi9zaApxZW11X2ZpbGU9Ii9ldGMvc3lzY29uZmlnL3FlbXUtZ2EiCmlmIFsgLWYgJHtxZW11X2ZpbGV9IF07IHRoZW4KICAgIHNlZCAtaSAtciAicy9eIz9CTEFDS0xJU1RfUlBDPS8jQkxBQ0tMSVNUX1JQQz0vIiAiJHtxZW11X2ZpbGV9IgogICAgaGFzX2dxYT0kKHN5c3RlbWN0bCBsaXN0LXVuaXRzIC0tZnVsbCAtYWxsIC10IHNlcnZpY2UgLS1wbGFpbiB8IGdyZXAgLW8gcWVtdS1ndWVzdC1hZ2VudC5zZXJ2aWNlKQogICAgaWYgW1sgLW4gJHtoYXNfZ3FhfSBdXTsgdGhlbgogICAgICAgIHN5c3RlbWN0bCByZXN0YXJ0IHFlbXUtZ3Vlc3QtYWdlbnQuc2VydmljZQogICAgZmkKZmkKCi0tPT09PT09PT09PT09PT09MjMwOTk4NDA1OTc0Mzc2MjQ3NT09LS0=", "host": "con01.vim1.local", "ephemeral_key_uuid": null, "architecture": null, "display_name": "dx_instance12", "system_metadata": {"image_os_distro": "centos", "clean_attempts": "3", "image_os_admin_user": "root", "image_image_type": "image", "image_os_version": "11", "boot_roles": "system_admin,courier_system_admin,ironic_system_admin,keystone_system_admin,nova_system_admin,octavia_system_admin,panko_system_admin,neutron_system_admin,glance_system_admin,cinder_system_admin,heat_system_admin,placement_system_admin", "owner_user_name": "dx", "image_hw_qemu_guest_agent": "yes", "image_min_disk": "0", "image_container_format": "bare", "image_min_ram": "0", "image_disk_format": "raw", "image_usage_type": "common", "owner_project_name": "dx_project"}, "task_state": null, "shutdown_terminate": false, "cell_name": null, "root_gb": 0, "locked": false, "name": "instance-00002ec1", "created_at": "2023-03-24T07:09:18.000000", "locked_by": null, "launch_index": 0, "memory_mb": 1024, "vcpus": 2, "image_ref": "a9cd271f-93b8-4e08-9a19-e328c46b1ac4", "root_device_name": "/dev/vda", "auto_disk_config": false, "os_type": null, "metadata": {}}}}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("attach_volume", consts.VOLUME, consts.CINDER, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.VOLUME || publisher.ErrorMsg == nil ||
		publisher.ResourceId != "74fbef54-f375-4c75-b058-8d41f850c3ec" ||
		publisher.ResourceState != consts.Available {
		t.Fatal("Failed to handleCinder")
	}
}

func TestHandleCinder_scheduler_create_volume(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"request_spec": {"backup_id": null, "snapshot_id": null, "volume_properties": {"status": "creating", "volume_type_id": "ae2d6670-757e-4400-9a24-e1c4548e9310", "group_id": null, "user_id": "a1e0d70fc2274d39af8434d633c3347e", "display_name": "dx_v1", "availability_zone": "nova", "reservations": ["7dbeb8ac-db74-45ba-b3e0-e6e555aa91fa", "f43755cf-30b0-4931-bbaf-048ea581f4ba", "ac46cccd-b5c0-4fd1-acbc-70067d4458b1", "bf99a428-65bf-465d-9049-0f29b041a4a8"], "multiattach": false, "attach_status": "detached", "source_volid": null, "cgsnapshot_id": null, "qos_specs": null, "encryption_key_id": null, "display_description": null, "snapshot_id": null, "consistencygroup_id": null, "project_id": "7e8babd4464e4c6da382a1a29d8da53a", "size": 10, "metadata": {}}, "source_volid": null, "cgsnapshot_id": null, "volume": {"migration_status": null, "provider_id": null, "availability_zone": "nova", "terminated_at": null, "updated_at": null, "provider_geometry": null, "replication_extended_status": null, "replication_status": null, "snapshot_id": null, "ec2_id": null, "deleted_at": null, "id": "4163962e-684d-4b37-92c2-7804dc8bc542", "size": 10, "display_name": "dx_v1", "display_description": null, "cluster_name": null, "metadata": {}, "name_id": "4163962e-684d-4b37-92c2-7804dc8bc542", "volume_admin_metadata": [], "encryption_key_id": null, "project_id": "7e8babd4464e4c6da382a1a29d8da53a", "launched_at": null, "scheduled_at": null, "status": "creating", "volume_type_id": "ae2d6670-757e-4400-9a24-e1c4548e9310", "multiattach": false, "deleted": false, "service_uuid": null, "provider_location": null, "volume_glance_metadata": [], "admin_metadata": {}, "host": null, "glance_metadata": {}, "consistencygroup_id": null, "source_volid": null, "provider_auth": null, "previous_status": null, "name": "volume-4163962e-684d-4b37-92c2-7804dc8bc542", "user_id": "a1e0d70fc2274d39af8434d633c3347e", "bootable": false, "shared_targets": true, "attach_status": "detached", "_name_id": null, "volume_metadata": [], "replication_driver_data": null, "group_id": null, "created_at": "2023-03-30T05:40:32.000000"}, "image_id": null, "availability_zones": ["nova"], "consistencygroup_id": null, "volume_type": {"name": "lvm", "qos_specs_id": null, "deleted": false, "created_at": "2022-04-28T08:27:14.000000", "updated_at": null, "extra_specs": {"volume_backend_name": "lvm"}, "is_public": true, "deleted_at": null, "id": "ae2d6670-757e-4400-9a24-e1c4548e9310", "projects": [], "description": null}, "volume_id": "4163962e-684d-4b37-92c2-7804dc8bc542", "operation": "create_volume", "resource_properties": {"status": "creating", "volume_type_id": "ae2d6670-757e-4400-9a24-e1c4548e9310", "group_id": null, "user_id": "a1e0d70fc2274d39af8434d633c3347e", "display_name": "dx_v1", "availability_zone": "nova", "reservations": ["7dbeb8ac-db74-45ba-b3e0-e6e555aa91fa", "f43755cf-30b0-4931-bbaf-048ea581f4ba", "ac46cccd-b5c0-4fd1-acbc-70067d4458b1", "bf99a428-65bf-465d-9049-0f29b041a4a8"], "multiattach": false, "attach_status": "detached", "source_volid": null, "cgsnapshot_id": null, "qos_specs": null, "encryption_key_id": null, "display_description": null, "snapshot_id": null, "consistencygroup_id": null, "project_id": "7e8babd4464e4c6da382a1a29d8da53a", "size": 10, "metadata": {}}, "group_id": null}, "volume_properties": {"status": "creating", "volume_type_id": "ae2d6670-757e-4400-9a24-e1c4548e9310", "group_id": null, "user_id": "a1e0d70fc2274d39af8434d633c3347e", "display_name": "dx_v1", "availability_zone": "nova", "reservations": ["7dbeb8ac-db74-45ba-b3e0-e6e555aa91fa", "f43755cf-30b0-4931-bbaf-048ea581f4ba", "ac46cccd-b5c0-4fd1-acbc-70067d4458b1", "bf99a428-65bf-465d-9049-0f29b041a4a8"], "multiattach": false, "attach_status": "detached", "source_volid": null, "cgsnapshot_id": null, "qos_specs": null, "encryption_key_id": null, "display_description": null, "snapshot_id": null, "consistencygroup_id": null, "project_id": "7e8babd4464e4c6da382a1a29d8da53a", "size": 10, "metadata": {}}, "reason": "NoValidBackend(u'No valid backend was found. No weighed backends available',)", "state": "error", "volume_id": "4163962e-684d-4b37-92c2-7804dc8bc542", "method": "create_volume"}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("scheduler.create_volume", consts.VOLUME, consts.CINDER, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.VOLUME || publisher.ErrorMsg == nil ||
		publisher.ResourceId != "4163962e-684d-4b37-92c2-7804dc8bc542" ||
		publisher.ResourceState != consts.Error {
		t.Fatal("Failed to handleCinder")
	}
}

func TestHandleCinder_create_snapshot_error(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"status": "error", "display_name": "dx_ss1", "availability_zone": "nova", "deleted": "", "tenant_id": "7e8babd4464e4c6da382a1a29d8da53a", "created_at": "2023-04-18T02:59:00+00:00", "snapshot_id": "b69f16f2-4a05-4a8c-a461-28ba1c0e94c0", "volume_type": "5ea726ad-c86c-4976-b854-bd7c67839ce3", "volume_size": 1, "volume_id": "f2d771f0-2959-42e8-a5ac-0b84407282c6", "user_id": "a1e0d70fc2274d39af8434d633c3347e", "error_msg": "Exception('create error',)", "metadata": ""}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("snapshot.create.error", consts.SNAPSHOT, consts.CINDER, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.SNAPSHOT || publisher.ErrorMsg == nil ||
		publisher.ResourceId != "b69f16f2-4a05-4a8c-a461-28ba1c0e94c0" ||
		publisher.ResourceState != consts.Error {
		t.Fatal("Failed to handleCinder")
	}
}

func TestHandleCinder_volume_resize_end(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"status": "in-use", "display_name": "dx_v3", "volume_attachment": [{"instance_uuid": "d61e0312-5ed4-4910-bd01-8a7d0692be62", "detach_time": null, "attach_time": "2023-04-23T08:26:31.000000", "deleted": false, "attach_mode": "rw", "created_at": "2023-04-23T08:26:22.000000", "attached_host": "con01.vim1.local", "updated_at": "2023-04-23T08:26:31.000000", "attach_status": "attached", "volume": {"migration_status": null, "provider_id": null, "availability_zone": "nova", "terminated_at": null, "updated_at": "2023-04-27T06:12:33.000000", "provider_geometry": null, "replication_extended_status": null, "replication_status": null, "snapshot_id": null, "ec2_id": null, "deleted_at": null, "id": "01c3011b-cc8c-42cb-9a36-95b0f5589695", "size": 15, "user_id": "d113fe136b044262b8e711f6fbecba15", "display_description": null, "cluster_name": null, "project_id": "e4a97e45239340b090fcebde8edcbfe5", "launched_at": "2023-04-14T07:59:32.000000", "scheduled_at": "2023-04-14T07:59:31.000000", "status": "in-use", "volume_type_id": "ef2e4472-b805-460b-a93f-cbc6a18ea0ac", "multiattach": false, "deleted": false, "service_uuid": "2145daf3-e558-4507-b7b2-d65c7305a41a", "provider_location": null, "host": "control@rbd-1#rbd-1", "consistencygroup_id": null, "source_volid": null, "provider_auth": null, "previous_status": "in-use", "display_name": "dx_v3", "bootable": false, "created_at": "2023-04-14T07:59:31.000000", "attach_status": "attached", "_name_id": null, "encryption_key_id": null, "replication_driver_data": null, "group_id": null, "shared_targets": false}, "connection_info": {"attachment_id": "283cddec-03c6-44ac-b8ba-003bbe93c1e6", "encrypted": false, "driver_volume_type": "rbd", "secret_uuid": "4eb5ed20-d3bd-4385-a7a1-0f632ab5d68e", "qos_specs": {"read_bytes_sec": "133120", "write_iops_sec": "5000", "write_bytes_sec": "133120", "read_iops_sec": "5000"}, "volume_id": "01c3011b-cc8c-42cb-9a36-95b0f5589695", "auth_username": "cinder", "secret_type": "ceph", "name": "volumes/volume-01c3011b-cc8c-42cb-9a36-95b0f5589695", "discard": true, "keyring": null, "cluster_name": "ceph", "auth_enabled": true, "hosts": ["10.50.31.1"], "access_mode": "rw", "ports": ["6789"]}, "volume_id": "01c3011b-cc8c-42cb-9a36-95b0f5589695", "mountpoint": "/dev/vdb", "deleted_at": null, "id": "283cddec-03c6-44ac-b8ba-003bbe93c1e6", "connector": {"initiator": "iqn.1994-05.com.redhat:5ef8ddcf4469", "ip": "10.50.31.1", "system uuid": "965D05B4-3383-03E2-11EB-3FA3E9333829", "platform": "x86_64", "host": "con01.vim1.local", "do_local_attach": false, "mountpoint": "/dev/vdb", "os_type": "linux2", "multipath": false}}], "availability_zone": "nova", "tenant_id": "e4a97e45239340b090fcebde8edcbfe5", "created_at": "2023-04-14T07:59:31+00:00", "volume_id": "01c3011b-cc8c-42cb-9a36-95b0f5589695", "volume_type": "ef2e4472-b805-460b-a93f-cbc6a18ea0ac", "host": "control@rbd-1#rbd-1", "replication_driver_data": null, "replication_status": null, "snapshot_id": null, "replication_extended_status": null, "user_id": "d113fe136b044262b8e711f6fbecba15", "metadata": [], "launched_at": "2023-04-14T07:59:32+00:00", "size": 15}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("volume.resize.end", consts.VOLUME, consts.CINDER, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.VOLUME || publisher.ErrorMsg != nil ||
		publisher.ResourceId != "01c3011b-cc8c-42cb-9a36-95b0f5589695" ||
		publisher.ResourceState != consts.InUse || publisher.VolumeSize != 15 {
		t.Fatal("Failed to handleCinder")
	}
}

func TestHandleCinder_detach_volume(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"exception": "{'message': u'Device detach failed for vdb: Unable to detach the device from the live config.', 'class': 'DeviceDetachFailed', 'kwargs': {'device': u'vdb', 'reason': u'Unable to detach the device from the live config.', 'code': 500}}", "args": {"volume_id": "01c3011b-cc8c-42cb-9a36-95b0f5589695"}}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("detach_volume", consts.VOLUME, consts.CINDER, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.VOLUME || publisher.ErrorMsg == nil ||
		publisher.ResourceId != "01c3011b-cc8c-42cb-9a36-95b0f5589695" ||
		publisher.ResourceState != consts.InUse {
		t.Fatal("Failed to handleCinder")
	}
}

func TestHandleCinder_volume_detach_end_for_instance_delete(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"status": "in-use", "display_name": "dx_v2", "volume_attachment": [], "availability_zone": "nova", "tenant_id": "e4a97e45239340b090fcebde8edcbfe5", "created_at": "2023-04-23T08:17:50+00:00", "volume_id": "9ece23a0-01b7-469f-b6ea-f0c83a14e747", "volume_type": "ef2e4472-b805-460b-a93f-cbc6a18ea0ac", "host": "control@rbd-1#rbd-1", "replication_driver_data": null, "replication_status": null, "snapshot_id": null, "replication_extended_status": null, "user_id": "d113fe136b044262b8e711f6fbecba15", "metadata": [], "launched_at": "2023-04-23T08:17:50+00:00", "size": 10}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("volume.detach.end", consts.VOLUME, consts.CINDER, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.VOLUME || publisher.ErrorMsg != nil ||
		publisher.ResourceId != "9ece23a0-01b7-469f-b6ea-f0c83a14e747" ||
		publisher.ResourceState != consts.Available {
		t.Fatal("Failed to handleCinder")
	}
}

func TestHandleCinder_volume_detach_failed_and_retry_end(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"status": "detaching", "display_name": "dx_v2", "volume_attachment": [], "availability_zone": "nova", "tenant_id": "e4a97e45239340b090fcebde8edcbfe5", "created_at": "2023-04-23T08:17:50+00:00", "volume_id": "9ece23a0-01b7-469f-b6ea-f0c83a14e747", "volume_type": "ef2e4472-b805-460b-a93f-cbc6a18ea0ac", "host": "control@rbd-1#rbd-1", "replication_driver_data": null, "replication_status": null, "snapshot_id": null, "replication_extended_status": null, "user_id": "d113fe136b044262b8e711f6fbecba15", "metadata": [], "launched_at": "2023-04-23T08:17:50+00:00", "size": 10}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("volume.detach.end", consts.VOLUME, consts.CINDER, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.VOLUME || publisher.ErrorMsg != nil ||
		publisher.ResourceId != "9ece23a0-01b7-469f-b6ea-f0c83a14e747" ||
		publisher.ResourceState != consts.Available {
		t.Fatal("Failed to handleCinder")
	}
}

// fip associate port
func TestHandleNeutron_sdn_floatingip_update_end(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"floatingip": {"router_id": "b409e1a2-f21c-493d-acb9-04f560422edf", "status": "ACTIVE", "qos_policy_binding": null, "floating_port_id": "ce3af58b-a18f-4e7c-b22e-aa5e494bdebd", "last_known_router_id": null, "fixed_port_id": "987eae63-fcfb-4210-9152-027031d25604", "floating_network_id": "e89828f4-21ca-4ebc-a624-7764895b3dc4", "standard_attr": {"description": "", "tags": [], "created_at": "2023-09-11T06:50:00.000000", "updated_at": "2023-09-12T06:02:56.000000", "revision_number": 1, "id": 11008, "resource_type": "floatingips"}, "fixed_ip_address": "192.197.197.23", "floating_ip_address": "10.50.98.181", "dns": null, "fixed_port": {"qos_network_policy_binding": null, "allowed_address_pairs": [], "distributed_port_binding": [], "device_owner": "compute:nova", "standard_attr_id": 11004, "fixed_ips": [{"subnet_id": "1637f9da-7954-4dd3-b40a-950e4c2af240", "network_id": "9327b327-f88a-4066-993e-0fe251ef318c", "port_id": "987eae63-fcfb-4210-9152-027031d25604", "ip_address": "192.197.197.23"}], "id": "987eae63-fcfb-4210-9152-027031d25604", "security_groups": [{"port_id": "987eae63-fcfb-4210-9152-027031d25604", "security_group_id": "8b054b8d-786f-451e-b2db-7fb916e198f7"}], "port_forwardings": [], "standard_attr": {"description": "", "tags": [], "created_at": "2023-09-11T06:39:20.000000", "updated_at": "2023-09-11T06:39:32.000000", "revision_number": 4, "id": 11004, "resource_type": "ports"}, "trunk_port": null, "mac_address": "fa:16:3e:ff:cc:fb", "sub_port": null, "project_id": "48ad435f0e8c44598d3236acdbb9ca47", "status": "ACTIVE", "binding_levels": [{"driver": "huawei_ac_ml2", "host": "con01.vim1.local", "segment_id": "9db1b852-cee3-41c8-9648-0fcece104060", "port_id": "987eae63-fcfb-4210-9152-027031d25604", "level": 0}, {"driver": "openvswitch", "host": "con01.vim1.local", "segment_id": "49b794e5-902f-4441-8c91-1fd01941ac5a", "port_id": "987eae63-fcfb-4210-9152-027031d25604", "level": 1}], "port_security": {"port_security_enabled": true, "port_id": "987eae63-fcfb-4210-9152-027031d25604"}, "qos_policy_binding": null, "device_id": "c809c113-b520-4a83-a789-829139510e4b", "port_bindings": [{"profile": "", "status": "ACTIVE", "vif_type": "ovs", "vif_details": "{\"datapath_type\": \"system\", \"ovs_hybrid_plug\": false, \"bridge_name\": \"br-int\", \"port_filter\": true, \"connectivity\": \"l2\"}", "vnic_type": "normal", "host": "con01.vim1.local", "port_id": "987eae63-fcfb-4210-9152-027031d25604"}], "name": "", "admin_state_up": true, "network_id": "9327b327-f88a-4066-993e-0fe251ef318c", "dns": null, "dhcp_opts": [], "ip_allocation": "immediate"}, "standard_attr_id": 11008, "project_id": "48ad435f0e8c44598d3236acdbb9ca47", "port_forwardings": [], "id": "a535a1de-ef48-41e6-9b06-550a340889ee"}}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("sdn.floatingip.update.end", consts.FLOATINGIP, consts.NEUTRON, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.FLOATINGIP || publisher.ErrorMsg != nil ||
		publisher.ResourceId != "a535a1de-ef48-41e6-9b06-550a340889ee" ||
		publisher.ResourceState != consts.Active {
		t.Fatal("Failed to handleNeutron")
	}
}

// fip disassociate port
func TestHandleNeutron_sdn_floatingip_update_end_disassociate(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"floatingip": {"router_id": null, "status": "DOWN", "description": "", "tags": [], "updated_at": "2023-09-12T06:03:30Z", "dns_domain": "", "floating_network_id": "e89828f4-21ca-4ebc-a624-7764895b3dc4", "port_forwardings": [], "fixed_ip_address": null, "floating_ip_address": "10.50.98.181", "revision_number": 2, "port_id": null, "id": "a535a1de-ef48-41e6-9b06-550a340889ee", "qos_policy_id": null, "tenant_id": "48ad435f0e8c44598d3236acdbb9ca47", "created_at": "2023-09-11T06:50:00Z", "port_details": null, "dns_name": "", "project_id": "48ad435f0e8c44598d3236acdbb9ca47"}}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("sdn.floatingip.update.end", consts.FLOATINGIP, consts.NEUTRON, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.FLOATINGIP || publisher.ErrorMsg != nil ||
		publisher.ResourceId != "a535a1de-ef48-41e6-9b06-550a340889ee" ||
		publisher.ResourceState != consts.Down {
		t.Fatal("Failed to handleNeutron")
	}
}

// fip port with qos policy
func TestHandleNeutron_port_update_end(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"port": {"allowed_address_pairs": [], "extra_dhcp_opts": [], "updated_at": "2023-09-12T06:21:32", "device_owner": "network:floatingip", "revision_number": 3, "binding:profile": {"fw_enabled": false}, "port_security_enabled": false, "fixed_ips": [{"subnet_id": "ebe38d71-08b5-4a2a-b00b-e55c0522fa1e", "ip_address": "10.50.98.222"}], "id": "79321696-9eac-470c-b4af-7ca82d5c2470", "security_groups": [], "binding:vif_details": {}, "binding:vif_type": "unbound", "mac_address": "fa:16:3e:4a:96:80", "project_id": "", "status": "ACTIVE", "binding:host_id": "", "description": "", "tags": [], "qos_policy_id": "760ff009-962c-49c6-8862-687267a653db", "name": "", "admin_state_up": true, "network_id": "e89828f4-21ca-4ebc-a624-7764895b3dc4", "tenant_id": "", "created_at": "2023-09-12T06:02:11Z", "binding:vnic_type": "normal", "device_id": "a98db965-0ad9-4cb2-bbc6-baee9dc4a9ce"}}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("port.update.end", consts.PORT, consts.NEUTRON, msgMap, true)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.FLOATINGIPPORT || publisher.ErrorMsg != nil ||
		publisher.ResourceId != "79321696-9eac-470c-b4af-7ca82d5c2470" ||
		publisher.ResourceState != consts.Active {
		t.Fatal("Failed to handleNeutron")
	}
}

// fip associate port in non-sdn env
func TestHandleNeutron_non_sdn_floatingip_update_end_associate(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"floatingip": {"router_id": "536002f6-7701-4821-923d-ebd6e49fbaa8", "status": "DOWN", "description": "", "tags": [], "port_id": "e618d843-3810-467a-9bf5-08006eca92ac", "created_at": "2023-11-01T12:43:45Z", "updated_at": "2023-11-02T01:34:31Z", "floating_network_id": "5ea7b8af-3cdf-4ad5-b353-516054b80ac9", "port_details": {"status": "ACTIVE", "name": "ak48", "admin_state_up": true, "network_id": "e1c17c9e-d496-4617-8b12-7d9871ad24a4", "device_owner": "compute:nova", "mac_address": "fa:16:3e:5b:bd:73", "device_id": "511098ca-a1f7-464e-8a7e-ef9612afb5de"}, "fixed_ip_address": "10.200.10.225", "floating_ip_address": "10.50.24.157", "revision_number": 7, "tenant_id": "b101be5c37ad4e98bc2543e15a61fce8", "project_id": "b101be5c37ad4e98bc2543e15a61fce8", "port_forwardings": [], "id": "d0743ed1-4de3-4c7b-b979-8517f68154d0", "qos_policy_id": "09a8a03e-1dff-4b1a-b848-d497ba153003"}}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("floatingip.update.end", consts.FLOATINGIP, consts.NEUTRON, msgMap, false)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.FLOATINGIP || publisher.ErrorMsg != nil ||
		publisher.ResourceId != "d0743ed1-4de3-4c7b-b979-8517f68154d0" ||
		publisher.ResourceState != consts.Active {
		t.Fatal("Failed to handleNeutron")
	}
}

// fip disassociate port in non-sdn env
func TestHandleNeutron_non_sdn_floatingip_update_end_disassociate(t *testing.T) {
	t.Helper()
	jsonBody := []byte(`{"payload": {"floatingip": {"router_id": null, "status": "ACTIVE", "description": "", "tags": [], "port_id": null, "created_at": "2023-11-01T12:43:45Z", "updated_at": "2023-11-02T01:35:32Z", "floating_network_id": "5ea7b8af-3cdf-4ad5-b353-516054b80ac9", "port_details": null, "fixed_ip_address": null, "floating_ip_address": "10.50.24.157", "revision_number": 8, "tenant_id": "b101be5c37ad4e98bc2543e15a61fce8", "project_id": "b101be5c37ad4e98bc2543e15a61fce8", "port_forwardings": [], "id": "d0743ed1-4de3-4c7b-b979-8517f68154d0", "qos_policy_id": "09a8a03e-1dff-4b1a-b848-d497ba153003"}}}`)
	msgMap := byteToMap(jsonBody)
	parser := NewParser("floatingip.update.end", consts.FLOATINGIP, consts.NEUTRON, msgMap, false)
	publisher := parser.Parse()
	fmt.Printf("%#v", publisher)
	if publisher.ResourceType != consts.FLOATINGIP || publisher.ErrorMsg != nil ||
		publisher.ResourceId != "d0743ed1-4de3-4c7b-b979-8517f68154d0" ||
		publisher.ResourceState != consts.Down {
		t.Fatal("Failed to handleNeutron")
	}
}
