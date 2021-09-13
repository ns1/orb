export const STRINGS = {
  header: {
    logotype: 'ORB',
  },
  footer: {
    disclaimer: '',
  },
  home: {
    title: 'Orb Observation Overview',
  },
  // Login Pages strings
  login: {
    presentation: {
      description: 'An Open-Source dynamic edge observability platform',
      action: 'Unleash the power of small data with dynamic edge observability',
    },
    form: {
      title: 'Log in or sign up',
      username: 'Username',
      password: 'Password',
      forgot: 'Forgot Password?',
      login: 'Log in',
      register: 'Register',
      nonuser: 'Don\'t have an account?',
    },
  },
  // Fleet Pages strings
  fleet: {
    title: 'Fleet Management',
  },
  // Sink Pages strings
  sink: {
    // sink statuses
    status: {
      active: 'Active',
      error: 'Error',
    },
    // sink.interface name descriptors
    propNames: {
      id: 'id',
      name: 'Name',
      description: 'Description',
      tags: 'Tags',
      status: 'Status',
      error: 'Error',
      backend: 'Service Type',
      config: 'Connection Details',
      config_remote_host: 'Remote Host',
      config_username: 'Username',
      config_password: 'Password',
      ts_created: 'Date Created',
    },
    // add page
    add: {
      header: 'Sink Creation',
    },
    // edit page
    edit: {
      header: 'Update Sink',
    },
    // delete modal
    delete: {
      header: 'Delete Sink Confirmation',
      body: 'Are you sure you want to delete this sink? This may cause Datasets which use this sink to become invalid. This action cannot be undone.',
      warning: '*To confirm, type your Sink name exactly as it appears',
      delete: 'I Understand, Delete This Sink',
      close: 'Close',
    },
    // details modal
    details: {
      header: 'Sink Details',
      close: 'Close',
    },
    // dashboard page
    list: {
      header: 'All Sinks',
      none: 'There are no sinks listed.',
      sink: 'sink',
      total: ['You have', 'total.'],
      error: 'have errors.',
      create: 'New Sink',
      filters: {
        select: 'Filter',
        name: 'Name',
        description: 'Description',
        status: 'Status',
        type: 'Type',
        tags: 'Tags',
      },
    },
  },
  // agents
  agentGroups: {
    // statuses
    status: {
      active: 'Active',
      error: 'Error',
    },
    // agent.interface name descriptors
    propNames: {
      id: 'id',
      name: 'Agent Group Name',
      description: 'Agent Group Description',
      key: 'Key',
      value: 'Value',
      tags: 'Tags',
      status: 'Status',
      error: 'Error',
      ts_created: 'Date Created',
      matches: 'Matches Against',
    },
    // matches
    match: {
      matchAny: 'The Selected Qualifiers Will Match Against',
      matchNone: 'The Selected Qualifiers Do Not Match Any Agent',
      agents: 'Agent(s)' +
        '',
      updated: 'Agent Group matches updated',
      expand: 'Expand',
      collapse: 'Collapse',
    },
    // add page
    add: {
      header: 'New Agent Group',
      step: {
        title1: 'Agent Group Details',
        desc1: 'This is how you will be able to easily identify your Agent Group',
        title2: 'Agent Group Tags',
        desc2: 'Set the tags that will be used to group Agents',
        title3: 'Review & Confirm',
      },
      success: 'Agent Group successfully created',
    },
    // edit page
    edit: {
      header: 'Update Agent',
    },
// delete modal
    delete: {
      header: 'Delete Agent Group Confirmation',
      body: 'Are you sure you want to delete this Agent Group? This may cause Datasets which use this Agent Group to become invalid. This action cannot be undone.',
      warning: '*To confirm, type the Agent Group name exactly as it appears',
      delete: 'I Understand, Delete This Agent Group',
      close: 'Close',
    },
    // details modal
    details: {
      header: 'Agent Group Details',
      close: 'Close',
    },
    // dashboard page
    list: {
      header: 'All Agents Groups',
      none: 'There are no agents listed.',
      agentGroup: 'agent',
      total: ['You have', 'total.'],
      error: 'have errors.',
      create: 'New Agent Group',
      filters: {
        select: 'Filter',
        name: 'Name',
        description: 'Description',
        status: 'Status',
        type: 'Type',
        tags: 'Tags',
      },
    },
  },
  // agent groups
  agents: {
    // statuses
    status: {
      active: 'Active',
      error: 'Error',
    },
    // agent.interface name descriptors
    propNames: {
      id: 'id',
      name: 'Agent Name',
      description: 'Agent Description',
      key: 'Key',
      value: 'Value',
      orb_tags: 'Orb Tags',
      state: 'Status',
      error: 'Error',
      ts_created: 'Date Created',
    },
    // add page
    add: {
      header: 'New Agent',
      step: {
        title1: 'Agent Details',
        desc1: 'This is how you will be able to easily identify your Agent',
        title2: 'Orb Tags',
        desc2: 'Set the tags that will be used to filter your agent',
        title3: 'Review & Confirm',
      },
      success: 'Agent successfully created',
    },
    // edit page
    edit: {
      header: 'Update Agent',
    },
// delete modal
    delete: {
      header: 'Delete Confirmation',
      body: 'Are you sure you want to delete this Agent? This action cannot be undone.',
      warning: '*To confirm, type the Agent label exactly as it appears',
      delete: 'I Understand, Delete This Agent',
      close: 'Close',
    },
    // details modal
    details: {
      header: 'Agent Details',
      close: 'Close',
    },
    // dashboard page
    list: {
      header: 'All Agents',
      none: 'There are no agents listed.',
      agentGroup: 'agent',
      total: ['You have', 'total.'],
      error: 'have errors.',
      create: 'New Agent',
      filters: {
        select: 'Filter',
        name: 'Name',
        description: 'Description',
        status: 'Status',
        type: 'Type',
        tags: 'Tags',
      },
    },
  },
  // stepper cues
  stepper: {
    back: 'Back',
    next: 'Next',
    save: 'Save',
  },
  // tags cues
  tags: {
    addTag: 'Add New Tag',
    key: 'Tag Key',
    value: 'Tag Value',
  },
};
