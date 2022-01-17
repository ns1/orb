import { Component } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { DynamicFormConfig } from 'app/common/interfaces/orb/dynamic.form.interface';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { PolicyTap } from 'app/common/interfaces/orb/policy/policy.tap.interface';

const CONFIG = {
  TAPS: 'TAPS',
  BACKEND: 'BACKEND',
  INPUTS: 'INPUTS',
  HANDLERS: 'HANDLERS',
  AGENT_POLICY: 'AGENT_POLICY',
};

@Component({
  selector: 'ngx-agent-policy-add-component',
  templateUrl: './agent.policy.add.component.html',
  styleUrls: ['./agent.policy.add.component.scss'],
})
export class AgentPolicyAddComponent {
  // #forms
  // agent policy general information - name, desc, backend
  detailsFG: FormGroup;

  // selected tap, input_type
  tapFG: FormGroup;

  // dynamic input config
  inputConfigFG: FormGroup;

  // dynamic input filter config
  inputFilterFG: FormGroup;

  // handlers
  handlerSelectorFG: FormGroup;

  dynamicHandlerConfigFG: FormGroup;
  dynamicHandlerFilterFG: FormGroup;

  // #key inputs holders
  // selected backend object
  backend: { [propName: string]: any };

  // selected tap object
  tap: PolicyTap;

  // selected input object
  input: {
    version?: string,
    config?: DynamicFormConfig,
    filter?: DynamicFormConfig,
  };

  // holds selected handler conf.
  // handler template currently selected, to be edited by user and then added to the handlers list or discarded
  liveHandler: {
    version?: string,
    config?: DynamicFormConfig,
    filter?: DynamicFormConfig,
    type?: string,
  };

  // holds all handlers added by user
  modules: {
    [propName: string]: {
      name?: string,
      type?: string,
      config?: { [propName: string]: {} | any },
      filter?: { [propName: string]: {} | any },
    },
  } = {};

  // #services responses
  // hold info retrieved
  availableBackends: {
    [propName: string]: {
      backend: string,
      description: string,
    },
  };

  availableTaps: { [propName: string]: PolicyTap };

  availableInputs: {
    [propName: string]: {
      version?: string,
      config?: DynamicFormConfig,
      filter?: DynamicFormConfig,
    },
  };

  availableHandlers: {
    [propName: string]: {
      version?: string,
      config?: DynamicFormConfig,
      filter?: DynamicFormConfig,
      metrics?: DynamicFormConfig,
      metrics_groups?: DynamicFormConfig,
    },
  } = {};

  // #if edit
  agentPolicy: AgentPolicy;

  agentPolicyID: string;

  isEdit: boolean;

  // #load controls
  isLoading = Object.entries(CONFIG)
    .reduce((acc, [value]) => {
      acc[value] = false;
      return acc;
    }, {});

  constructor(
    private agentPoliciesService: AgentPoliciesService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
  ) {
    this.agentPolicyID = this.route.snapshot.paramMap.get('id');
    this.agentPolicy = this.newAgent();
    this.isEdit = !!this.agentPolicyID;

    this.readyForms();

    Promise.all([
      this.isEdit ? this.retrieveAgentPolicy() : Promise.resolve(),
      this.getBackendsList(),
    ]).catch(reason => console.warn(`Couldn't fetch data. Reason: ${ reason }`))
      .then(() => this.updateForms())
      .catch((reason) => console.warn(`Couldn't fetch ${ this.agentPolicy?.backend } data. Reason: ${ reason }`));
  }

  newAgent() {
    return {
      name: '',
      description: '',
      backend: 'pktvisor',
      tags: {},
      version: 1,
      policy: {
        kind: 'collection',
        input: {
          config: {},
          tap: '',
          input_type: '',
        },
        handlers: {
          modules: {},
        },
      },
    } as AgentPolicy;
  }

  retrieveAgentPolicy() {
    return new Promise(resolve => {
      this.agentPoliciesService.getAgentPolicyById(this.agentPolicyID).subscribe(policy => {
        this.agentPolicy = policy;
        this.isLoading[CONFIG.AGENT_POLICY] = false;
        resolve(policy);
      });
    });
  }

  readyForms() {
    const {
      name,
      description,
      backend,
      policy: {
        input: {
          tap,
          input_type,
        },
        handlers: {
          modules,
        },
      },
    } = this.agentPolicy;

    this.modules = modules;

    this.detailsFG = this._formBuilder.group({
      name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$')]],
      description: [description],
      backend: [{ value: backend, disabled: backend !== '' }, [Validators.required]],
    });
    this.tapFG = this._formBuilder.group({
      selected_tap: [tap, Validators.required],
      input_type: [input_type, Validators.required],
    });

    this.handlerSelectorFG = this._formBuilder.group({
      'selected_handler': ['', [Validators.required]],
      'label': ['', [Validators.required]],
    });

    this.dynamicHandlerConfigFG = this._formBuilder.group({});
    this.dynamicHandlerFilterFG = this._formBuilder.group({});
  }

  updateForms() {
    const {
      name,
      description,
      backend,
      policy: {
        handlers: {
          modules,
        },
      },
    } = this.agentPolicy;

    this.detailsFG.patchValue({ name, description, backend });

    this.modules = modules;

    this.dynamicHandlerConfigFG = this._formBuilder.group({});
    this.dynamicHandlerFilterFG = this._formBuilder.group({});

    this.onBackendSelected(backend).catch(reason => console.warn(`${ reason }`));


  }

  getBackendsList() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.BACKEND] = true;
      this.agentPoliciesService.getAvailableBackends().subscribe(backends => {
        this.availableBackends = !!backends && backends.reduce((acc, curr) => {
          acc[curr.backend] = curr;
          return acc;
        }, {});

        this.isLoading[CONFIG.BACKEND] = false;

        resolve(backends);
      });
    });
  }

  onBackendSelected(selectedBackend) {
    return new Promise((resolve) => {
      this.backend = this.availableBackends[selectedBackend];
      this.backend['config'] = {};

      // todo hardcoded for pktvisor
      this.getBackendData().then(() => {
        resolve(null);
      });
    });
  }

  getBackendData() {
    return Promise.all([this.getTaps(), this.getInputs(), this.getHandlers()])
      .then(value => {
        if (this.isEdit && this.agentPolicy) {
          const selected_tap = this.agentPolicy.policy.input.tap;
          this.tapFG.patchValue({ selected_tap }, { emitEvent: true });
          this.onTapSelected(selected_tap);
          this.tapFG.controls.selected_tap.disable();
        }

      }, reason => console.warn(`Cannot retrieve backend data - reason: ${ JSON.parse(reason) }`))
      .catch(reason => {
        console.warn(`Cannot retrieve backend data - reason: ${ JSON.parse(reason) }`);
      });
  }

  getTaps() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.TAPS] = true;
      this.agentPoliciesService.getBackendConfig([this.backend.backend, 'taps'])
        .subscribe(taps => {
          this.availableTaps = taps.reduce((acc, curr) => {
            acc[curr.name] = curr;
            return acc;
          }, {});

          this.isLoading[CONFIG.TAPS] = false;

          resolve(taps);
        });
    });
  }

  onTapSelected(selectedTap) {
    this.tap = this.availableTaps[selectedTap];
    this.tapFG.controls.selected_tap.patchValue(selectedTap);

    const { input } = this.agentPolicy.policy;
    const { input_type, config_predefined, filter_predefined } = this.tap;

    this.tap.config = {
      ...config_predefined,
      ...input.config,
    };

    this.tap.filter = {
      ...filter_predefined,
      ...input.filter,
    };

    if (input_type) {
      this.onInputSelected(input_type);
    } else {
      this.input = null;
      this.tapFG.controls.input_type.reset('');
    }
  }

  getInputs() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.INPUTS] = true;
      this.agentPoliciesService.getBackendConfig([this.backend.backend, 'inputs'])
        .subscribe(inputs => {
          this.availableInputs = !!inputs && inputs;

          this.isLoading[CONFIG.INPUTS] = false;

          resolve(inputs);
        });
    });

  }

  onInputSelected(input_type) {
    // TODO version here
    this.input = this.availableInputs[input_type]['1.0'];

    this.tapFG.patchValue({ input_type });

    // input type config model
    const { config: inputConfig, filter: filterConfig } = this.input;
    // if editing, some values might not be overrideable any longer, all should be prefilled in form
    const { config: agentConfig, filter: agentFilter } = this.agentPolicy.policy.input;
    // tap config values, cannot be overridden if set
    const {config_predefined: preConfig, filter_predefined: preFilter} = this.tap;

    // populate form controls for config
    const inputConfDynamicCtrl = Object.entries(inputConfig)
      .reduce((acc, [key, input]) => {
        const value = agentConfig?.[key] || '';
        if (!preConfig?.includes(key)) {
          acc[key] = [
            value,
            [!!input?.props?.required && input.props.required === true ? Validators.required : Validators.nullValidator],
          ];
        }
        return acc;
      }, {});

    this.inputConfigFG = Object.keys(inputConfDynamicCtrl).length > 0 ? this._formBuilder.group(inputConfDynamicCtrl) : null;

    const inputFilterDynamicCtrl = Object.entries(filterConfig)
      .reduce((acc, [key, input]) => {
        const value = !!agentFilter?.[key] ? agentFilter[key] : '';
        // const disabled = !!preConfig?.[key];
        if (!preFilter?.includes(key)) {
          acc[key] = [
            value,
            [!!input?.props?.required && input.props.required === true ? Validators.required : Validators.nullValidator],
          ];
        }
        return acc;
      }, {});

    this.inputFilterFG = Object.keys(inputFilterDynamicCtrl).length > 0 ? this._formBuilder.group(inputFilterDynamicCtrl) : null;

  }

  getHandlers() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.HANDLERS] = true;

      this.agentPoliciesService.getBackendConfig([this.backend.backend, 'handlers'])
        .subscribe(handlers => {
          this.availableHandlers = handlers || {};

          this.isLoading[CONFIG.HANDLERS] = false;
          resolve(handlers);
        });
    });
  }

  onHandlerSelected(selectedHandler) {
    if (this.dynamicHandlerConfigFG) {
      this.dynamicHandlerConfigFG = null;
    }
    if (this.dynamicHandlerFilterFG) {
      this.dynamicHandlerFilterFG = null;
    }

    // TODO - hardcoded for v: 1.0 -: always retrieve latest
    this.liveHandler = selectedHandler !== '' && !!this.availableHandlers[selectedHandler] ?
      { ...this.availableHandlers[selectedHandler]['1.0'], type: selectedHandler } : null;

    const { config, filter } = this.liveHandler || { config: {}, filter: {} };

    const dynamicConfigControls = Object.entries(config || {}).reduce((controls, [key]) => {
      controls[key] = ['', [Validators.required]];
      return controls;
    }, {});

    this.dynamicHandlerConfigFG = Object.keys(dynamicConfigControls).length > 0 ? this._formBuilder.group(dynamicConfigControls) : null;

    const dynamicFilterControls = Object.entries(filter || {}).reduce((controls, [key]) => {
      controls[key] = ['', [Validators.required]];
      return controls;
    }, {});

    const suggestName = this.getSeriesHandlerName(selectedHandler);
    this.handlerSelectorFG.patchValue({label: suggestName});

    this.dynamicHandlerFilterFG = Object.keys(dynamicFilterControls).length > 0 ? this._formBuilder.group(dynamicFilterControls) : null;
  }

  getSeriesHandlerName(handlerType) {
    const count = 1 + Object.values(this.modules || {}).filter(({type}) => type === handlerType).length;
    return `handler_${handlerType}_${count}`;
  }

  checkValidName() {
    const { value } = this.handlerSelectorFG.controls.label;
    const hasTagForKey = Object.keys(this.modules).find(key => key === value);
    return value && value !== '' && !hasTagForKey;
  }

  onHandlerAdded() {
    let config = {};
    let filter = {};

    if (this.dynamicHandlerConfigFG !== null) {
      config = Object.entries(this.dynamicHandlerConfigFG.controls)
        .reduce((acc, [key, control]) => {
          if (control.value) acc[key] = control.value;
          return acc;
        }, {});
    }

    if (this.dynamicHandlerFilterFG !== null) {
      filter = Object.entries(this.dynamicHandlerFilterFG.controls)
        .reduce((acc, [key, control]) => {
          if (control.value) acc[key] = control.value;
          return acc;
        }, {});
    }

    const handlerName = this.handlerSelectorFG.controls.label.value;
    this.modules[handlerName] = ({
      type: this.liveHandler.type,
      config,
      filter,
    });

    this.handlerSelectorFG.reset({
      selected_handler: { value: '', disabled: false },
      label: { value: '', disabled: false },
    });
    this.onHandlerSelected('');
  }

  onHandlerRemoved(handlerName) {
    delete this.modules[handlerName];
  }

  goBack() {
    this.router.navigateByUrl('/pages/datasets/policies');
  }

  onFormSubmit() {
    const payload = {
      name: this.detailsFG.controls.name.value,
      description: this.detailsFG.controls.description.value,
      backend: this.detailsFG.controls.backend.value,
      tags: {},
      version: !!this.isEdit && !!this.agentPolicy.version && this.agentPolicy.version || 1,
      policy: {
        kind: 'collection',
        input: {
          tap: this.tap.name,
          input_type: this.tapFG.controls.input_type.value,
          ...Object.entries(this.inputConfigFG.controls)
            .map(([key, control]) => ({ [key]: control.value }))
            .reduce((acc, curr) => {
              for (const [key, value] of Object.entries(curr)) {
                if (!!value && value !== '') acc.config[key] = value;
              }
              return acc;
            }, {config: {}}),
          ...Object.entries(this.inputFilterFG.controls)
            .map(([key, control]) => ({ [key]: control.value }))
            .reduce((acc, curr) => {
              for (const [key, value] of Object.entries(curr)) {
                if (!!value && value !== '') acc.filter[key] = value;
              }
              return acc;
            }, {filter: {}}),
        },
        handlers: {
          modules: Object.entries(this.modules).reduce((acc, [key, value]) => {
            const {type, config, filter} = value;
            acc[key] = {
              type: type,
              config: Object.entries(config).length > 0 && config || undefined,
              filter: Object.entries(filter).length > 0 && filter || undefined,
            };
            if (Object.keys(config || {}).length > 0) acc[key][config] = config;
            return acc;
          }, {}),
        },
      },
    } as AgentPolicy;

    if (Object.keys(payload.policy?.input?.config).length <= 0)
      delete payload.policy.input.config;
    if (Object.keys(payload.policy?.input?.filter).length <= 0)
      delete payload.policy.input.filter;

    if (this.isEdit) {
      // updating existing sink
      this.agentPoliciesService.editAgentPolicy({ ...payload, id: this.agentPolicyID }).subscribe(() => {
        this.notificationsService.success('Agent Policy successfully updated', '');
        this.goBack();
      });
    } else {
      this.agentPoliciesService.addAgentPolicy(payload).subscribe(() => {
        this.notificationsService.success('Agent Policy successfully created', '');
        this.goBack();
      });
    }
  }
}
