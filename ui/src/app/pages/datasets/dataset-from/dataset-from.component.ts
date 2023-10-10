import { ChangeDetectorRef, Component, Input, OnChanges, OnInit, SimpleChange, SimpleChanges } from '@angular/core';
import {
  AbstractControl,
  FormBuilder,
  FormGroup,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { NbDialogRef, NbDialogService } from '@nebular/theme';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { DatasetDeleteComponent } from 'app/pages/datasets/delete/dataset.delete.component';
import { AgentMatchComponent } from 'app/pages/fleet/agents/match/agent.match.component';
import { Observable, of } from 'rxjs';

export const DATASET_RESPONSE = {
  EDITED: 'edited',
  CANCELED: 'canceled',
  DELETED: 'deleted',
  CREATED: 'created',
};

const CONFIG = {
  SINKS: 'SINKS',
  GROUPS: 'GROUPS',
  POLICIES: 'POLICIES',
  DATASET: 'DATASET',
};

@Component({
  selector: 'ngx-dataset-from',
  templateUrl: './dataset-from.component.html',
  styleUrls: ['./dataset-from.component.scss'],
})
export class DatasetFromComponent implements OnInit, OnChanges {
  @Input()
  dataset: Dataset;

  @Input()
  policy: AgentPolicy;

  @Input()
  group: AgentGroup;

  isEdit: boolean;

  isGroupSelected: boolean = false;

  selectedGroup: string;
  groupName: string;
  selectedPolicy: string;
  fetchedData: boolean;
  sinkIDs: string[];
  availableAgentGroups: AgentGroup[];
  filteredAgentGroups$: Observable<AgentGroup[]>;
  availableAgentPolicies: AgentPolicy[];
  availableSinks: Sink[];
  unselectedSinks: Sink[];
  form: FormGroup;
  loading = Object.entries(CONFIG).reduce((acc, [value]) => {
    acc[value] = false;
    return acc;
  }, {});

  isRequesting: boolean;

  constructor(
    private agentGroupsService: AgentGroupsService,
    private agentPoliciesService: AgentPoliciesService,
    private datasetService: DatasetPoliciesService,
    private sinksService: SinksService,
    private notificationsService: NotificationsService,
    private fb: FormBuilder,
    private dialogRef: NbDialogRef<DatasetFromComponent>,
    private dialogService: NbDialogService,
    private cdr: ChangeDetectorRef,
  ) {
    this.isEdit = false;
    this.groupName = '';
    this.availableAgentGroups = [];
    this.fetchedData = false;
    this.filteredAgentGroups$ = of(this.availableAgentGroups);
    this.availableAgentPolicies = [];
    this.availableSinks = [];
    this._selectedSinks = [];
    this.unselectedSinks = [];
    this.sinkIDs = [];
    this.isRequesting = false;

    this.getDatasetAvailableConfigList();

    this.readyForms();

    this.form.get('agent_group_id').valueChanges.subscribe(value => {
      this.ngOnChanges({ agent_group_id: new SimpleChange(null, value, true) });
    });
  }

  private _selectedSinks: Sink[];

  get selectedSinks(): Sink[] {
    return this._selectedSinks;
  }

  // #load controls

  set selectedSinks(sinks: Sink[]) {
    this._selectedSinks = sinks;
    this.sinkIDs = sinks.map((sink) => sink.id);
    this.form.controls.sink_ids.patchValue(this.sinkIDs);
    this.form.controls.sink_ids.markAsDirty();
    this.updateUnselectedSinks();
  }

  readyForms() {
    const { agent_policy_id, agent_group_id, sink_ids } =
      this?.dataset ||
      ({
        agent_group_id: '',
        agent_policy_id: '',
        sink_ids: [],
      } as Dataset);

    this.form = this.fb.group({
      agent_policy_id: [agent_policy_id, [Validators.required]],
      agent_group_id: [agent_group_id, [Validators.required]],
      agent_group_name: [null, [this.groupNameValidator]],
      sink_ids: [sink_ids],
    });
  }

  groupNameValidator = (): ValidatorFn => {
    return (control: AbstractControl) =>
      this.availableAgentGroups.filter((agent) => agent.name === control.value)
        .length === 0
        ? { noMatch: 'Select a valid agent' }
        : null;
  }

  updateFormSelectedAgentGroupId(groupName) {
    const group = this.availableAgentGroups.filter(
      (agent) => agent.name === groupName,
    );
    let id;
    if (group.length > 0) {
      id = group[0].id;
    }
    this.form.patchValue({ agent_group_id: id });
    this.cdr.markForCheck();
  }

  updateFormSelectedAgentGroupName(groupId) {
    const group = this.availableAgentGroups.filter(
      (agent) => agent.id === groupId,
    );
    if (group.length > 0) {
      this.groupName = group[0].name;
    }
    this.form.patchValue({ agent_group_name: this.groupName });
    this.cdr.markForCheck();
  }

  onChangeGroupName(event) {
    const value = event.currentTarget.value;
    this.onFilterGroup(value);
  }

  onSelectChangeGroupName(event) {
    this.onFilterGroup(event);
  }

  onFilterGroup(value) {
    this.filteredAgentGroups$ = of(this.filter(value));
  }

  onMatchingAgentsModal() {
    this.dialogService.open(AgentMatchComponent, {
      context: {
        agentGroupId: this.form.controls.agent_group_id.value,
        policy: this.policy,
      },
      autoFocus: true,
      closeOnEsc: true,
    });
  }
  ngOnChanges(changes: SimpleChanges): void {
    if (changes.agent_group_id.currentValue) {
      this.isGroupSelected = true;
    } else {
      this.isGroupSelected = false;
    }
  }

  ngOnInit(): void {
    if (!!this.group) {
      this.selectedGroup = this.group.id;
      this.form.patchValue({ agent_group_id: this.group.id });
      this.form.controls.agent_group_id.disable();
    }
    if (!!this.policy) {
      this.selectedPolicy = this.policy.id;
      this.form.patchValue({ agent_policy_id: this.policy.id });
      this.form.controls.agent_policy_id.disable();
    }
    if (!!this.dataset) {
      const { name, agent_group_id, agent_policy_id, sink_ids } = this.dataset;
      this.selectedGroup = agent_group_id;
      this.selectedSinks =
        (!!sink_ids &&
          this.availableSinks.filter((sink) => sink_ids.includes(sink.id))) ||
        [];
      this.selectedPolicy = agent_policy_id;
      this.form.patchValue({ name, agent_group_id, agent_policy_id, sink_ids });
      this.isEdit = true;
      this.form.controls.agent_group_id.disable();
      this.form.controls.agent_policy_id.disable();

      this.unselectedSinks = this.availableSinks.filter(
        (sink) => !this._selectedSinks.includes(sink),
      );
    }
  }

  updateUnselectedSinks() {
    this.unselectedSinks = this.availableSinks.filter(
      (sink) => !this._selectedSinks.includes(sink),
    );
  }

  getDatasetAvailableConfigList() {
    Promise.all([
      this.getAvailableAgentGroups(),
      this.getAvailableAgentPolicies(),
      this.getAvailableSinks(),
    ])
      .then(
        (value) => {
          this.fetchedData = true;
        },
        (reason) =>
          console.warn(
            `Cannot retrieve available configurations - reason: ${JSON.parse(
              reason,
            )}`,
          ),
      )
      .catch((reason) => {
        console.warn(
          `Cannot retrieve backend data - reason: ${JSON.parse(reason)}`,
        );
      });
  }

  getAvailableAgentGroups() {
    return new Promise((resolve) => {
      this.loading[CONFIG.GROUPS] = true;
      this.agentGroupsService
        .getAllAgentGroups()
        .subscribe((resp: AgentGroup[]) => {
          this.availableAgentGroups = resp.sort((a, b) =>
            a.name > b.name ? -1 : 1,
          );
          this.filteredAgentGroups$ = of(this.availableAgentGroups);
          this.loading[CONFIG.GROUPS] = false;

          if (this.dataset?.agent_group_id) {
            this.updateFormSelectedAgentGroupName(this.dataset.agent_group_id);
          }
          resolve(this.availableAgentGroups);
        });
    });
  }

  getAvailableAgentPolicies() {
    return new Promise((resolve) => {
      this.loading[CONFIG.POLICIES] = true;

      this.agentPoliciesService
        .getAllAgentPolicies()
        .subscribe((resp: AgentPolicy[]) => {
          this.availableAgentPolicies = resp;
          this.loading[CONFIG.POLICIES] = false;

          resolve(this.availableAgentPolicies);
        });
    });
  }

  getAvailableSinks() {
    return new Promise((resolve) => {
      this.loading[CONFIG.SINKS] = true;
      this.sinksService.getAllSinks().subscribe((resp: Sink[]) => {
        this.availableSinks = resp;
        const selectedSinkIds = this.dataset?.sink_ids || [];
        this.selectedSinks = selectedSinkIds.map((sink) => {
          return resp.find((anotherSink) => anotherSink.id === sink);
        });

        this.loading[CONFIG.SINKS] = false;

        resolve(this.availableSinks);
      });
    });
  }

  isLoading() {
    return Object.values<boolean>(this.loading).reduce(
      (prev, curr) => prev && curr,
    );
  }

  onFormSubmit() {
    this.isRequesting = true;
    const payload = {
      name: this.createNewName(),
      agent_group_id: this.form.controls.agent_group_id.value,
      agent_policy_id: this.form.controls.agent_policy_id.value,
      sink_ids: this._selectedSinks.map((sink) => sink.id),
    } as Dataset;
    if (this.isEdit) {
      // updating existing dataset
      this.datasetService
        .editDataset({ ...payload, id: this.dataset.id })
        .subscribe(() => {
          this.notificationsService.success('Dataset successfully updated', '');
          this.dialogRef.close(DATASET_RESPONSE.EDITED);
        });
    } else {
      this.datasetService.addDataset(payload).subscribe(() => {
        this.notificationsService.success('Dataset successfully created', '');
        this.dialogRef.close(DATASET_RESPONSE.CREATED);
      });
    }
  }

  onDelete() {
    this.dialogService
      .open(DatasetDeleteComponent, {
        context: { name: this.dataset.name },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.datasetService.deleteDataset(this.dataset.id).subscribe(() => {
            this.notificationsService.success(
              'Dataset successfully deleted',
              '',
            );
            this.dialogRef.close(DATASET_RESPONSE.DELETED);
          });
        }
      });
  }

  onClose() {
    this.dialogRef.close(DATASET_RESPONSE.CANCELED);
  }

  private filter(value: string): AgentGroup[] {
    let filtered;
    if (value === '') {
      filtered = this.availableAgentGroups;
    } else {
      filtered = this.availableAgentGroups.filter((group) =>
        group.name.includes(value),
      );
    }
    this.updateFormSelectedAgentGroupId(value);

    return filtered;
  }

  createNewName() {
    const ts = Date.now();
    return `dataset_${ts}`;
  }
}
