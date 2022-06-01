// This file was auto-generated from a YAML file

package v5_0_5

func init() {
	Profile["/cloudify/5.0.5/profile.yaml"] = `
# Content below is modified from:
# https://cloudify.co/spec/cloudify/5.0.5/types.yaml

imports:
- data.yaml

##################################################################################
# The types defined here make use of features only available in the
# 'cloudify_dsl_1_3' tosca definitions version
##################################################################################
tosca_definitions_version: cloudify_dsl_1_3

##################################################################################
# Cloudify basic types metadata
##################################################################################
metadata:
  cloudify_types: true
  puccini.scriptlet.import:tosca.lib.utils: internal:/tosca/common/1.0/js/lib/utils.js
  puccini.scriptlet.import:tosca.lib.traversal: internal:/tosca/common/1.0/js/lib/traversal.js
  puccini.scriptlet.import:tosca.coerce: internal:/tosca/common/1.0/js/coerce.js
  puccini.scriptlet.import:tosca.resolve: internal:/tosca/common/1.0/js/resolve.js

##################################################################################
# Base type definitions
##################################################################################
node_types:

  # base type for provided cloudify types
  cloudify.nodes.Root:
    interfaces:
      cloudify.interfaces.lifecycle:
        precreate: {}
        create: {}
        configure: {}
        start: {}
        poststart: {}
        prestop: {}
        stop: {}
        delete: {}
        postdelete: {}
      cloudify.interfaces.validation:
        create: {}
        delete: {}
        # Deprecated - should not be implemented.
        creation: {}
        deletion: {}
      cloudify.interfaces.monitoring:
        start: {}
        stop: {}

  # A host (physical / virtual or LXC) in a topology
  cloudify.nodes.Compute:
    derived_from: cloudify.nodes.Root
    properties:
      ip:
        default: ''
      os_family:
        description: |
          Property specifying what type of operating system family
          this compute node will run.
        default: linux
      agent_config:
        type: cloudify.datatypes.AgentConfig
        default:
          install_method: remote
          port: 22
          network: default

      # DEPRECATED
      install_agent:
        default: ''
      cloudify_agent:
        default: {}

    interfaces:
      cloudify.interfaces.cloudify_agent:

        #####################################################
        # This operation will download the necessary files
        # for the agent and place them were needed.
        # ---------------------------------------------------
        # It will be executed on the manager and will connect
        # to the agent via fabric/winrm
        #####################################################

        create:
          implementation: agent.cloudify_agent.installer.operations.create
          executor: central_deployment_agent

        #####################################################
        # This operation will create all the necessary
        # configuration for the agent to be started.
        # ---------------------------------------------------
        # It will be executed on the manager and will connect
        # to the agent via fabric/winrm
        #####################################################

        configure:
          implementation: agent.cloudify_agent.installer.operations.configure
          executor: central_deployment_agent

        #####################################################
        # This operation will actually start the agent
        # process.
        # ---------------------------------------------------
        # It will be executed on the manager and will connect
        # to the agent via fabric/winrm
        #####################################################

        start:
          implementation: agent.cloudify_agent.installer.operations.start
          executor: central_deployment_agent

        #####################################################
        # This operation will stop the agent process.
        # ---------------------------------------------------
        # It will be executed on the manager and will connect
        # to the agent via fabric/winrm
        #####################################################

        stop:
          implementation: agent.cloudify_agent.installer.operations.stop
          executor: central_deployment_agent

        #####################################################
        # This operation will stop the agent process
        # without connecting to the remote host.
        # ---------------------------------------------------
        # It will be executed on the agent host via AMQP
        #####################################################

        stop_amqp:
          implementation: agent.cloudify_agent.operations.stop
          executor: host_agent

        #####################################################
        # This operation will delete agent resources.
        # ---------------------------------------------------
        # It will be executed on the manager and will connect
        # to the agent via fabric/winrm
        #####################################################

        delete:
          implementation: agent.cloudify_agent.installer.operations.delete
          executor: central_deployment_agent

        #####################################################
        # This operation will restart the agent process.
        # ---------------------------------------------------
        # It will be executed on the manager and will connect
        # to the agent via fabric/winrm
        #####################################################

        restart:
          implementation: agent.cloudify_agent.installer.operations.restart
          executor: central_deployment_agent

        #####################################################
        # This operation will restart the agent process
        # without connecting to the remote host.
        # the agent will have a new name after this
        # operation, but will of course listen to the same
        # queue.
        # ---------------------------------------------------
        # It will be executed on the agent host via AMQP
        #####################################################

        restart_amqp:
          implementation: agent.cloudify_agent.operations.restart
          executor: host_agent

        #####################################################
        # This operation will install additional plugins
        # on the agent environment.
        # ---------------------------------------------------
        # It will be executed directly on the agent via AMQP.
        # not requiring any new remote connection to be
        # established.
        #####################################################

        install_plugins:
          implementation: agent.cloudify_agent.operations.install_plugins
          executor: host_agent

        #####################################################
        # This operation will uninstall plugins
        # on the agent environment.
        # ---------------------------------------------------
        # It will be executed directly on the agent via AMQP.
        # not requiring any new remote connection to be
        # established.
        #####################################################

        uninstall_plugins:
          implementation: agent.cloudify_agent.operations.uninstall_plugins
          executor: host_agent

        #####################################################
        # This operation will install new agent on specific
        # machine using existing agent.
        # ---------------------------------------------------
        # It will be executed on central deployment agent.
        # This agent will send installation request via AMQP
        # to an agent that was previously installed.
        # This request might be sent to other manager's
        # rabbitmq server.
        #####################################################

        create_amqp:
          implementation: agent.cloudify_agent.operations.create_agent_amqp
          executor: central_deployment_agent
          inputs:
            install_agent_timeout:
              default: 300
            manager_ip:
              description: The private ip of the new manager
              default: ''
            manager_certificate:
              description: The cloudify_internal_ca_cert.pem of the new manager
              default: ''
            stop_old_agent:
              description: Stop the old agent after installing the new one
              default: false


        #####################################################
        # This operation attempts to communicate with an
        # agent that has already been installed. This
        # operation should not have side effects.
        # ---------------------------------------------------
        # It will be executed on the central deployment
        # agent. This request might be sent to a different
        # manager's RabbitMQ server
        #####################################################

        validate_amqp:
          implementation: agent.cloudify_agent.operations.validate_agent_amqp
          executor: central_deployment_agent
          inputs:
            validate_agent_timeout:
              default: 20

      cloudify.interfaces.host:  # DEPRECATED
        get_state: {}

      cloudify.interfaces.monitoring_agent:
        install: {}
        start: {}
        stop: {}
        uninstall: {}

  # A Linux container with or without docker
  cloudify.nodes.Container:
    derived_from: cloudify.nodes.Compute

  # A tier in a topology
  cloudify.nodes.Tier:
    derived_from: cloudify.nodes.Root

  # A storage volume in a topology
  cloudify.nodes.Volume:
    derived_from: cloudify.nodes.Root

  # A file system a volume should be formatted to
  cloudify.nodes.FileSystem:
    derived_from: cloudify.nodes.Root
    properties:
      use_external_resource:
        description: >
          Enables the use of already formatted volumes.
        type: boolean
        default: false
      partition_type:
        description: >
          The partition type. 83 is a Linux Native Partition.
        type: integer
        default: 83
      fs_type:
        description: >
          The type of the File System.
          Supported types are [ext2, ext3, ext4, fat, ntfs, swap]
        type: string
      fs_mount_path:
        description: >
          The path of the mount point.
        type: string
    interfaces:
      cloudify.interfaces.lifecycle:
        configure:
          implementation: script.script_runner.tasks.run
          inputs:
            script_path:
              default: file:///opt/manager/resources/cloudify/fs/mkfs.sh

  # A storage Container (Object Store segment)
  cloudify.nodes.ObjectStorage:
    derived_from: cloudify.nodes.Root

  # An isolated virtual layer 2 domain or a logical / virtual switch
  cloudify.nodes.Network:
    derived_from: cloudify.nodes.Root

  # An isolated virtual layer 3 subnet with IP range
  cloudify.nodes.Subnet:
    derived_from: cloudify.nodes.Root

  cloudify.nodes.Port:
    derived_from: cloudify.nodes.Root

  # A network router
  cloudify.nodes.Router:
    derived_from: cloudify.nodes.Root

  # A virtual Load Balancer
  cloudify.nodes.LoadBalancer:
    derived_from: cloudify.nodes.Root

  # A virtual floating IP
  cloudify.nodes.VirtualIP:
    derived_from: cloudify.nodes.Root

  # A security group
  cloudify.nodes.SecurityGroup:
    derived_from: cloudify.nodes.Root

  # A middleware component in a topology
  cloudify.nodes.SoftwareComponent:
    derived_from: cloudify.nodes.Root

  cloudify.nodes.DBMS:
    derived_from: cloudify.nodes.SoftwareComponent

  cloudify.nodes.Database:
    derived_from: cloudify.nodes.Root

  cloudify.nodes.WebServer:
    derived_from: cloudify.nodes.SoftwareComponent
    properties:
      port:
        default: 80

  cloudify.nodes.ApplicationServer:
    derived_from: cloudify.nodes.SoftwareComponent

  cloudify.nodes.MessageBusServer:
    derived_from: cloudify.nodes.SoftwareComponent

  # An application artifact to deploy
  cloudify.nodes.ApplicationModule:
    derived_from: cloudify.nodes.Root

  # A type for a Cloudify Manager, to be used in manager blueprints
  cloudify.nodes.CloudifyManager:
    derived_from: cloudify.nodes.SoftwareComponent
    properties:
      cloudify:
        description: >
          Configuration for Cloudify Manager
        default:
          resources_prefix: ''
          cloudify_agent:
            min_workers: 2
            max_workers: 5
            remote_execution_port: 22
            user: ubuntu
          workflows:
            task_retries: -1  # this means forever
            task_retry_interval: 30
          policy_engine:
            start_timeout: 30
      cloudify_packages:
        description: >
          Links to Cloudify packages to be installed on the manager

  cloudify.nodes.Component:
    derived_from: cloudify.nodes.Root
    properties:
      resource_config:
        type: cloudify.datatypes.Component
        default: {}
      client:
        description: >
          Cloudify HTTP client configuration,
          if empty the current Cloudify manager client will be used.
        type: cloudify.datatypes.RemoteCloudifyManagerClient
        default: {}
      plugins:
        description: >
          Dictionary of plugins to upload,
          which each plugin is in format of:
            plugin-name:
              wagon_path: Url for plugin wagon file,
              plugin_yaml_path: Url for plugin yaml file
        type: dict
        default: {}
      secrets:
        description: >
          Dictionary of secrets to set before deploying Components,
          which each secret is in format of:
            secret-name: value
        type: dict
        default: {}
    interfaces:
      cloudify.interfaces.lifecycle:
        create:
          implementation: cfy_extensions.cloudify_types.component.upload_blueprint
        configure:
          implementation: cfy_extensions.cloudify_types.component.create
        start:
          implementation: cfy_extensions.cloudify_types.component.execute_start
          inputs:
            workflow_id:
              type: string
              default: install
            timeout:
              description: How long (in seconds) to wait for execution to finish before timing out.
              type: integer
              default: 1800
            interval:
              description: Polling interval (seconds).
              type: integer
              default: 10
        stop:
          implementation: cfy_extensions.cloudify_types.component.execute_start
          inputs:
            workflow_id:
              default: uninstall
            resource_config:
              default:
                blueprint: { get_property: [ SELF, resource_config, blueprint ] }
                deployment: { get_property: [ SELF, resource_config, deployment ] }
                executions_start_args:
                  allow_custom_parameters: true
        delete:
          implementation: cfy_extensions.cloudify_types.component.delete

  cloudify.nodes.ServiceComponent:
    derived_from: cloudify.nodes.Component

  cloudify.nodes.SharedResource:
    derived_from: cloudify.nodes.Root
    properties:
      resource_config:
        type: cloudify.datatypes.SharedResource
        default: {}
      client:
        description: >
          Client configuration, if empty Cloudify manager client will be used.
        type: cloudify.datatypes.RemoteCloudifyManagerClient
        default: {}
    interfaces:
      cloudify.interfaces.lifecycle:
        create:
          implementation: cfy_extensions.cloudify_types.shared_resource.connect_deployment


##################################################################################
# Base relationship definitions
##################################################################################
relationships:

  cloudify.relationships.depends_on:
    source_interfaces:
      cloudify.interfaces.relationship_lifecycle:
        preconfigure: {}
        postconfigure: {}
        establish: {}
        unlink: {}
    target_interfaces:
      cloudify.interfaces.relationship_lifecycle:
        preconfigure: {}
        postconfigure: {}
        establish: {}
        unlink: {}
    properties:
      connection_type:
        default: all_to_all

  cloudify.relationships.connected_to:
    derived_from: cloudify.relationships.depends_on

  cloudify.relationships.contained_in:
    derived_from: cloudify.relationships.depends_on

  #####################################################
  # This relationship will create a dependency of that
  # node's lifecycle to a specific operation in the
  # lifecycle of the specified target node.
  #####################################################
  cloudify.relationships.depends_on_lifecycle_operation:
    derived_from: cloudify.relationships.depends_on
    properties:
      operation:
        description: The target's lifecycle operation name.
        type: string

  cloudify.relationships.depends_on_shared_resource:
    derived_from: cloudify.relationships.depends_on
    target_interfaces:
      cloudify.interfaces.relationship_lifecycle:
        establish:
          implementation: cfy_extensions.cloudify_types.shared_resource.execute_workflow
          inputs:
            workflow_id:
              description: >
                The workflow id that will be run in the SharedResource's deployment as
                implementation defined there.
              type: string
            parameters:
              description: >
                Inputs for running the workflow in the format of key-value dictionary.
              type: dict
              default: {}
            timeout:
              description: >
                Timeout in seconds for running the specified workflow on the deployment.
              type: integer
              default: 10
        unlink:
          implementation: cfy_extensions.cloudify_types.shared_resource.execute_workflow
          inputs:
            workflow_id:
              description: >
                The workflow id that will be run in the SharedResource's deployment as
                implementation defined there.
              type: string
            parameters:
              description: >
                Inputs for running the workflow in the format of key-value dictionary.
              type: dict
              default: {}
            timeout:
              description: >
                Timeout in seconds for running the specified workflow on the deployment.
              type: integer
              default: 10

  cloudify.relationships.connected_to_shared_resource:
    derived_from: cloudify.relationships.connected_to
    target_interfaces:
      cloudify.interfaces.relationship_lifecycle:
        establish:
          implementation: cfy_extensions.cloudify_types.shared_resource.execute_workflow
          inputs:
            workflow_id:
              description: >
                The workflow id that will be run in the SharedResource's deployment as
                implementation defined there.
              type: string
            parameters:
              description: >
                Inputs for running the workflow in the format of key-value dictionary.
              type: dict
              default: {}
            timeout:
              description: >
                Timeout in seconds for running the specified workflow on the deployment.
              type: integer
              default: 10
        unlink:
          implementation: cfy_extensions.cloudify_types.shared_resource.execute_workflow
          inputs:
            workflow_id:
              description: >
                The workflow id that will be run in the SharedResource deployment as
                defined there.
              type: string
            parameters:
              description: >
                Inputs for running the workflow in the format of key-value dictionary.
              type: dict
              default: {}
            timeout:
              description: >
                Timeout in seconds for running the specified workflow on the deployment.
              type: integer
              default: 10

  cloudify.relationships.file_system_depends_on_volume:
    derived_from: cloudify.relationships.depends_on
    source_interfaces:
      cloudify.interfaces.relationship_lifecycle:
        preconfigure:
          implementation: script.script_runner.tasks.run
          inputs:
            script_path:
              default: file:///opt/manager/resources/cloudify/fs/fdisk.sh
            device_name:
              default: { get_attribute: [TARGET, device_name] }

  cloudify.relationships.file_system_contained_in_compute:
    derived_from: cloudify.relationships.contained_in
    source_interfaces:
      cloudify.interfaces.relationship_lifecycle:
        establish:
          implementation: script.script_runner.tasks.run
          inputs:
            script_path:
              default: file:///opt/manager/resources/cloudify/fs/mount.sh
        unlink:
          implementation: script.script_runner.tasks.run
          inputs:
            script_path:
              default: file:///opt/manager/resources/cloudify/fs/unmount.sh

##################################################################################
# Workflows
##################################################################################
workflows:

  install:
    mapping: default_workflows.cloudify.plugins.workflows.install
    is_cascading: false

  update:
    mapping: default_workflows.cloudify.plugins.workflows.update
    is_cascading: false
    parameters:
      update_id:
        default: ''
      skip_install:
        default: false
      skip_uninstall:
        default: false
      added_instance_ids:
        default: []
        type: list
      added_target_instances_ids:
        default: []
        type: list
      removed_instance_ids:
        default: []
        type: list
      remove_target_instance_ids:
        default: []
        type: list
      modified_entity_ids:
        default: []
        type: list
      extended_instance_ids:
        default: []
        type: list
      extend_target_instance_ids:
        default: []
        type: list
      reduced_instance_ids:
        default: []
        type: list
      reduce_target_instance_ids:
        default: []
        type: list
      ignore_failure:
        default: false
        type: boolean
      install_first:
        default: false
        type: boolean
      node_instances_to_reinstall:
        default: []
        type: list
      central_plugins_to_install:
        default: []
        type: list
      central_plugins_to_uninstall:
        default: []
        type: list
      update_plugins:
        default: true
        type: boolean

  uninstall:
    mapping: default_workflows.cloudify.plugins.workflows.uninstall
    is_cascading: false
    parameters:
      ignore_failure:
        default: false
        type: boolean

  start:
    mapping: default_workflows.cloudify.plugins.workflows.start
    is_cascading: false
    parameters:
      operation_parms:
        default: {}
      run_by_dependency_order:
        default: true
      type_names:
        default: []
      node_ids:
        default: []
      node_instance_ids:
        default: []

  stop:
    mapping: default_workflows.cloudify.plugins.workflows.stop
    is_cascading: false
    parameters:
      operation_parms:
        default: {}
      run_by_dependency_order:
        default: true
      type_names:
        default: []
      node_ids:
        default: []
      node_instance_ids:
        default: []

  restart:
    mapping: default_workflows.cloudify.plugins.workflows.restart
    is_cascading: false
    parameters:
      stop_parms:
        default: {}
      start_parms:
        default: {}
      run_by_dependency_order:
        default: true
      type_names:
        default: []
      node_ids:
        default: []
      node_instance_ids:
        default: []

  execute_operation:
    mapping: default_workflows.cloudify.plugins.workflows.execute_operation
    is_cascading: false
    parameters:
      operation: {}
      operation_kwargs:
        default: {}
      allow_kwargs_override:
        default: null
      run_by_dependency_order:
        default: false
      type_names:
        default: []
      node_ids:
        default: []
      node_instance_ids:
        default: []

  heal:
    mapping: default_workflows.cloudify.plugins.workflows.auto_heal_reinstall_node_subgraph
    is_cascading: false
    parameters:
      node_instance_id:
        description: Which node instance has failed
      diagnose_value:
        description: Diagnosed reason of failure
        default: Not provided
      ignore_failure:
        default: true
        type: boolean

  scale:
    mapping: default_workflows.cloudify.plugins.workflows.scale_entity
    is_cascading: false
    parameters:
      scalable_entity_name:
        description: >
          Which node/group to scale. Note that the parameter specified
          should denote the node/group name and NOT the node/group instance id.
      delta:
        description: >
            How many node/group instances should be added/removed.
            A positive number denotes increase of instances.
            A negative number denotes decrease of instances.
        default: 1
        type: integer
      scale_compute:
        description: >
            If a node name is passed as the ` + "`" + `scalable_entity_name` + "`" + ` parameter
            and that node is contained (transitively) within a compute node
            and this property is 'true', operate on the compute node instead
            of the specified node.
        default: false
      include_instances:
        description: >
            A node instance ID or list of node instance IDs to prioritise for
            scaling down.
            If a larger amount of included instances are provided than the delta
            then one of the included instances will be scaled down while the
            others will be ignored.
            If a smaller amount of included instances are provided than the delta
            then the remaining instances will be selected arbitrarily.
            This is only valid for scaling down and cannot be used with
            scale_compute set to true.
        default: null
      exclude_instances:
        description: >
            A node instance ID or list of node instance IDs which must not be
            removed when scaling down.
            If the amount of excluded instances plus the absolute delta is equal
            to or greater than the total amount of instances then the scaling
            operation will fail and no nodes will be scaled down.
            This is only valid for scaling down and cannot be used with
            scale_compute set to true.
        default: null
      ignore_failure:
        default: false
        type: boolean
      rollback_if_failed:
        description: >
            If this is False then no rollback will be triggered when an error
            occurs during the workflow, otherwise the rollback will be
            triggered.
        default: true
        type: boolean

  install_new_agents:
    mapping: default_workflows.cloudify.plugins.workflows.install_new_agents
    is_cascading: false
    parameters:
      install_agent_timeout:
        default: 300
      node_ids:
        default: []
      node_instance_ids:
        default: []
      install_methods:
        default: null
      validate:
        default: True
      install:
        default: True
      install_script:
        default: ''
      manager_ip:
        description: The private ip of the new manager
        default: ''
      manager_certificate:
        description: The cloudify_internal_ca_cert.pem of the new manager
        default: ''
      stop_old_agent:
        description: Stop the old agent after the new agent is installed
        default: false

  validate_agents:
    mapping: default_workflows.cloudify.plugins.workflows.validate_agents
    is_cascading: false
    parameters:
      node_ids:
        default: []
      node_instance_ids:
        default: []
      install_methods:
        default: null

##################################################################################
# Base data type definitions
##################################################################################
data_types:

  cloudify.datatypes.AgentConfig:
    description: >
      Cloudify agent configuration schema.
    properties:
      install_method:
        description: |
          Specifies how (and if) the cloudify agent should be installed.
          Valid values are:
          * none - No agent will be installed on the host.
          * remote - An agent will be installed using SSH on linux hosts and WinRM on windows hosts.
          * init_script - An agent will be installed via a script that will run on the host when it gets created.
                          This method is only supported for specific IaaS plugins.
          * plugin - An agent will be installed via a plugin which will run a script on the host.
                     This method is only supported for specific IaaS plugins.
          * provided - An agent is assumed to already be installed on the host image.
                       That agent will be configured and started via a script that will run on the host when it gets created.
                       This method is only supported for specific IaaS plugins.
        type: string
        required: true
      service_name:
        description: |
          Used to set the the cloudify agent service name.

          If not set, the default value for the service name is:
          - Linux: 'cloudify-worker-<id>'
          - Windows: '<id>'

          where 'id' is the instance id of the compute node in which the agent is running.

          Note: the value in this field, takes precedence over the deprecated
          'cloudify.nodes.Compute.cloudify_agent.name'.
        type: string
        required: false
      network:
        description: >
          The name of the manager network to which the agent should be
          connected. By default, the value will be ` + "`" + `default` + "`" + ` (which is the
          manager's private IP, by default)
        type: string
        required: false
      user:
        description: >
          For host agents, the agent will be installed for this user.
        type: string
        required: false
      key:
        description: >
          For host agents that are installed via SSH, this is the path to the private
          key that will be used to connect to the host.
          In most cases, this value will be derived automatically during bootstrap.
        type: string
        required: false
      password:
        description: >
          For host agents that are installed via SSH (on linux) and WinRM (on windows)
          this property can be used to connect to the host.
          For linux hosts, this property is optional in case the key property is properly configured
          (either explicitly or implicitly during bootstrap).
          For windows hosts that are installed via WinRM, this property is also optional
          and depends on whether the password runtime property has been set by the relevant IaaS plugin,
          prior to the agent installation.
        type: string
        required: false
      port:
        description: >
          For host agents that are installed via SSH (on linux) and WinRM (on windows),
          this is the port used to connect to the host.
          The default values are 22 for linux hosts and 5985 for windows hosts.
        type: integer
        required: false
      process_management:
        description: >
          Process management specific configuration. (type: dictionary)
        required: false
      min_workers:
        description: >
          Minimum number of agent workers. By default, the value will be 0.
          Note: For windows based agents, this property is ignored and min_workers is set to the value of max_workers.
        type: integer
        required: false
      max_workers:
        description: >
          Maximum number of agent workers. By default, the value will be 5.
        type: integer
        required: false
      heartbeat:
        description: >
          The interval of the AMQP heartbeats in seconds
        required: false
      disable_requiretty:
        description: >
          For linux based agents, disables the requiretty setting in the sudoers file. By default, this value will be true.
        type: boolean
        required: false
      env:
        description: >
          Optional environment variables that the agent will be started with. (type: dictionary)
        required: false
      extra:
        description: >
          Optional additional low level configuration details. (type: dictionary)
        required: false
      executable_temp_path:
        description: >
          Directory to use for temporary executable files. This is useful for installations where the
          default temporary directory (` + "`" + `/tmp` + "`" + `) is mounted with ` + "`" + `noexec` + "`" + `.
        type: string
        required: false
      log_level:
        description: >
          The logging level for the agent. Can be one of the following values: critical, error, warning,
          info, debug
        type: string
        required: false
      log_max_bytes:
        description: >
          Maximum number of bytes in the agent's log file, before it is rolled over.
        type: integer
        required: false
      log_max_history:
        description: >
          Maximum number of historical log files to keep.
        type: integer
        required: false

  cloudify.datatypes.Blueprint:
    properties:
      external_resource:
        description: >
          Use external blueprint resource.
        type: boolean
        default: false
      id:
        description: >
          This is the blueprint ID that the Component's node is connected to.
        type: string
        required: false
      main_file_name:
        description: >
          The application blueprint filename. If the blueprint consists many
          imported files this is the main blueprint.
        type: string
        default: blueprint.yaml
      blueprint_archive:
        description: >
          The URL of a .zip to upload to the manager
          (Will be skipped if external_resource == True).
        type: string
        default: ""
        required: true

  cloudify.datatypes.ComponentDeployment:
    properties:
      id:
        description: >
          This is the deployment ID that the Component's node is connected to.
        type: string
        required: false
      inputs:
        description: >
          The inputs to the deployment.
        type: dict
        default: {}
      logs:
        description: >
          This is a flag for logs and events redirect from the deployment, by default true.
        type: boolean
        required: false
      auto_inc_suffix:
        description: >
          Optional, will add a suffix to the given deployment ID in the form of an auto incremented index.
        type: boolean
        required: false

  cloudify.datatypes.Component:
    properties:
      blueprint:
        type: cloudify.datatypes.Blueprint
        required: true
      deployment:
        type: cloudify.datatypes.ComponentDeployment
        required: true
      executions_start_args:
        description: >
          Optional params for Component executions.
        type: dict
        default: {}

  cloudify.datatypes.RemoteCloudifyManagerClient:
    properties:
      host:
        description: >
          Host of Cloudify's manager machine.
        type: string
        required: false
      port:
        description: >
          The port of REST API service on Cloudify's management machine.
        type: integer
        required: false
      protocol:
        description: >
          The protocol of REST API service on management machine, defaults to http.
        type: string
        required: false
      api_version:
        description: >
          The version of Cloudify REST API service.
        type: string
        required: false
      headers:
        description: >
          Headers to be added to HTTP requests.
        type: dict
        required: false
      query_params:
        description: >
          Query parameters to be added to the HTTP request.
        type: dict
        required: false
      cert:
        description: >
          Path on the Cloudify manager to a copy of the target Cloudify manager's certificate.
        type: string
        required: false
      trust_all:
        description: >
          If False, the server's certificate (self-signed or not) will be verified.
        type: boolean
        required: false
      username:
        description: >
          Cloudify user username.
        type: string
        required: false
      password:
        description: >
          Cloudify user password.
        type: string
        required: false
      token:
        description: >
          Cloudify user token.
        type: string
        required: false
      tenant:
        description: >
          Cloudify user accessible tenant name.
        type: string
        required: false

  cloudify.datatypes.SharedResourceDeployment:
    properties:
      id:
        description: >
          This is the deployment ID that the SharedResource node is connected to.
        type: string
        required: true

  cloudify.datatypes.SharedResource:
    properties:
      deployment:
        type: cloudify.datatypes.SharedResourceDeployment
        required: true

##################################################################################
# Base artifact definitions
##################################################################################
plugins:
  cfy_extensions:
    executor: central_deployment_agent
    install: false

  agent:
    executor: central_deployment_agent
    install: false

  default_workflows:
    executor: central_deployment_agent
    install: false

  script:
    executor: host_agent
    install: false

##################################################################################
# Policy types definitions
##################################################################################
policy_types:

    cloudify.policies.types.host_failure:
        properties: &BASIC_AH_POLICY_PROPERTIES
            policy_operates_on_group:
                description: |
                    If the policy should maintain its state for the whole group
                    or each node instance individually.
                default: false
            is_node_started_before_workflow:
                description: Before triggering workflow, check if the node state is started
                default: true
            interval_between_workflows:
                description: |
                    Trigger workflow only if the last workflow was triggered earlier than interval-between-workflows seconds ago.
                    if < 0  workflows can run concurrently.
                default: 300
            service:
                description: Service names whose events should be taken into consideration
                default:
                    - service
        source: file:///opt/manager/resources/cloudify/policies/host_failure.clj

    cloudify.policies.types.threshold:
        properties: &THRESHOLD_BASED_POLICY_PROPERTIES
            <<: *BASIC_AH_POLICY_PROPERTIES
            service:
                description: The service name
                default: service
            threshold:
                description: The metric threshold value
            upper_bound:
                description: |
                    boolean value for describing the semantics of the threshold.
                    if 'true': metrics whose value is bigger than the threshold will cause the triggers to be processed.
                    if 'false': metrics with values lower than the threshold will do so.
                default: true
            stability_time:
                description: How long a threshold must be breached before the triggers will be processed
                default: 0
        source: file:///opt/manager/resources/cloudify/policies/threshold.clj

    cloudify.policies.types.ewma_stabilized:
        properties:
            <<: *THRESHOLD_BASED_POLICY_PROPERTIES
            ewma_timeless_r:
                description: |
                    r is the ratio between successive events. The smaller it is, the smaller impact on the computed value the most recent event has.
                default: 0.5
        source: file:///opt/manager/resources/cloudify/policies/ewma_stabilized.clj

##################################################################################
# Policy triggers definitions
##################################################################################
policy_triggers:

  cloudify.policies.triggers.execute_workflow:
    parameters:
      workflow:
        description: Workflow name to execute
      workflow_parameters:
        description: Workflow paramters
        default: {}
      force:
        description: |
          Should the workflow be executed even when another execution
          for the same workflow is currently in progress
        default: false
      allow_custom_parameters:
        description: |
          Should parameters not defined in the workflow parameters
          schema be accepted
        default: false
      socket_timeout:
        description: Socket timeout when making request to manager REST in ms
        default: 1000
      conn_timeout:
        description: Connection timeout when making request to manager REST in ms
        default: 1000
    source: file:///opt/manager/resources/cloudify/triggers/execute_workflow.clj
`
}
