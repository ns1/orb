from hamcrest import *
import requests
from control_plane_agents import base_orb_url
from behave import given, when, then
from utils import random_string, filter_list_by_parameter_start_with

policy_name_prefix = "test_policy_name_"
default_handler = "net"
handle_label = "default_" + default_handler


@when("a new policy is created")
def create_new_policy(context):
    context.policy_name = policy_name_prefix + random_string(10)
    context.policy = create_policy(context.token, context.policy_name, handle_label, default_handler)



@then("referred policy must be listened on the orb policies list")
def check_policies(context):
    policy_id = context.policy['id']
    policy = get_policy(context.token, policy_id)
    assert_that(policy['name'], equal_to(context.policy_name), "Incorrect policy name")


@then('cleanup policies')
def clean_policies(context):
    """
    Remove all policies starting with 'policy_name_prefix' from the orb

    :param context: Behave class that contains contextual information during the running of tests.
    """
    token = context.token
    policies_list = list_policies(token)
    policies_filtered_list = filter_list_by_parameter_start_with(policies_list, 'name', policy_name_prefix)
    delete_policies(token, policies_filtered_list)


@given("that a policy already exists")
def new_policy(context):
    create_new_policy(context)
    check_policies(context)


def create_policy(token, policy_name, handler_label, handler, description=None, tap="default_pcap",
                  input_type="pcap", host_specification=None, filter_expression=None, backend_type="pktvisor"):
    """

    Creates a new policy in Orb control plane


    :param (str) token: used for API authentication
    :param (str) policy_name:  of the policy to be created
    :param (str) handler_label:  of the handler
    :param (str) handler: to be added
    :param (str) description: description of policy
    :param tap: named, host specific connection specifications for the raw input streams accessed by pktvisor
    :param input_type: this must reference a tap name, or application of the policy will fail
    :param (str) host_specification: Subnets (comma separated) which should be considered belonging to this host,
    in CIDR form. Used for ingress/egress determination, defaults to host attached to the network interface.
    :param filter_expression: these decide exactly which data to summarize and expose for collection
    :param backend_type: Agent backend this policy is for. Cannot change once created. Default: pktvisor
    :return: (dict) a dictionary containing the created policy data
    """
    json_request = {"name": policy_name, "description": description, "backend": backend_type,
                    "policy": {"kind": "collection", "input": {"tap": tap, "input_type": input_type},
                               "handlers": {"modules": {handler_label: {"type": handler}}}},
                    "config": {"host_spec": host_specification}, "filter": {"bpf": filter_expression}}
    headers_request = {'Content-type': 'application/json', 'Accept': '*/*', 'Authorization': token}

    response = requests.post(base_orb_url + '/api/v1/policies/agent', json=json_request, headers=headers_request)
    assert_that(response.status_code, equal_to(201),
                'Request to create policy failed with status=' + str(response.status_code))

    return response.json()


def get_policy(token, policy_id):
    """
    Gets a policy from Orb control plane

    :param (str) token: used for API authentication
    :param (str) policy_id: that identifies policy to be fetched
    :returns: (dict) the fetched policy
    """

    get_policy_response = requests.get(base_orb_url + '/api/v1/policies/agent/' + policy_id,
                                       headers={'Authorization': token})

    assert_that(get_policy_response.status_code, equal_to(200),
                'Request to get policy id=' + policy_id + ' failed with status=' + str(get_policy_response.status_code))

    return get_policy_response.json()


def list_policies(token, limit=100):
    """
    Lists all policies from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :returns: (list) a list of policies
    """
    response = requests.get(base_orb_url + '/api/v1/policies/agent', headers={'Authorization': token},
                            params={'limit': limit})

    assert_that(response.status_code, equal_to(200),
                'Request to list policies failed with status=' + str(response.status_code))

    policies_as_json = response.json()
    return policies_as_json['data']


def delete_policies(token, list_of_policies):
    """
    Deletes from Orb control plane the policies specified on the given list

    :param (str) token: used for API authentication
    :param (list) list_of_policies: that will be deleted
    """

    for policy in list_of_policies:
        delete_policy(token, policy['id'])


def delete_policy(token, policy_id):
    """
    Deletes a policy from Orb control plane

    :param (str) token: used for API authentication
    :param (str) policy_id: that identifies the policy to be deleted
    """

    response = requests.delete(base_orb_url + '/api/v1/policies/agent/' + policy_id,
                               headers={'Authorization': token})

    assert_that(response.status_code, equal_to(204), 'Request to delete policy id='
                + policy_id + ' failed with status=' + str(response.status_code))

