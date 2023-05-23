import {
  ChangeDetectorRef,
  Component,
  OnChanges,
  OnDestroy,
  OnInit,
  ViewChild,
} from '@angular/core';
import {
  ActivatedRoute,
  NavigationEnd,
  Router,
  RouterEvent,
} from '@angular/router';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { PolicyConfig } from 'app/common/interfaces/orb/policy/config/policy.config.interface';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { OrbService } from 'app/common/services/orb.service';
import { PolicyDetailsComponent } from 'app/shared/components/orb/policy/policy-details/policy-details.component';
import { PolicyInterfaceComponent } from 'app/shared/components/orb/policy/policy-interface/policy-interface.component';
import { STRINGS } from 'assets/text/strings';
import { Subscription } from 'rxjs';
import yaml from 'js-yaml';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { filter } from 'rxjs/operators';

@Component({
  selector: 'ngx-agent-view',
  templateUrl: './agent.policy.view.component.html',
  styleUrls: ['./agent.policy.view.component.scss'],
})
export class AgentPolicyViewComponent implements OnInit, OnDestroy, OnChanges {
  strings = STRINGS.agents;

  isLoading: boolean;

  policyId: string;

  policy: AgentPolicy;

  datasets: Dataset[];
  groups: AgentGroup[];

  policySubscription: Subscription;

  editMode = {
    details: false,
    interface: false,
  };

  @ViewChild(PolicyDetailsComponent) detailsComponent: PolicyDetailsComponent;

  @ViewChild(PolicyInterfaceComponent)
  interfaceComponent: PolicyInterfaceComponent;

  constructor(
    private route: ActivatedRoute,
    private policiesService: AgentPoliciesService,
    private orb: OrbService,
    private cdr: ChangeDetectorRef,
    private notifications: NotificationsService,
    private router: Router,
  ) {}

  ngOnInit() {
    this.fetchData();
  }

  fetchData() {
    this.isLoading = true;
    this.policyId = this.route.snapshot.paramMap.get('id');
    this.retrievePolicy();
  }

  ngOnChanges(): void {
    this.fetchData();
  }

  isEditMode() {
    return Object.values(this.editMode).reduce(
      (prev, cur) => prev || cur,
      false,
    );
  }

  canSave() {
    const detailsValid = this.editMode.details
      ? this.detailsComponent?.formGroup?.status === 'VALID'
      : true;

    const interfaceValid = this.editMode.interface
      ? this.interfaceComponent?.formControl?.status === 'VALID'
      : true;

    return detailsValid && interfaceValid;
  }

  discard() {
    this.editMode.details = false;
    this.editMode.interface = false;
  }

  save() {
    const { format, version, name, description, id, backend } = this.policy;

    // get values from all modified sections' forms and submit through service.
    const policyDetails = this.detailsComponent.formGroup?.value;
    const tags = this.detailsComponent.selectedTags;
    const policyInterface = this.interfaceComponent.code;

    // trying to work around rest api
    const detailsPartial = (!!this.editMode.details && {
      ...policyDetails,
    }) || { name, description };

    let interfacePartial = {};

    try {
      if (format === 'yaml') {
        yaml.load(policyInterface);

        interfacePartial = {
          format,
          policy_data: policyInterface,
        };
      } else {
        interfacePartial = {
          policy: JSON.parse(policyInterface) as PolicyConfig,
        };
      }

      const payload = {
        ...detailsPartial,
        ...interfacePartial,
        version,
        id,
        tags,
        backend,
      } as AgentPolicy;

      this.policiesService.editAgentPolicy(payload).subscribe((resp) => {
        this.discard();
        this.policy = resp;
        this.fetchData();
      });

      this.notifications.success('Agent Policy updated successfully', '');
    } catch (err) {
      this.notifications.error(
        'Failed to edit Agent Policy',
        `Error: Invalid ${format.toUpperCase()}`,
      );
    }
  }

  retrievePolicy() {
    this.policySubscription = this.orb
      .getPolicyFullView(this.policyId)
      .subscribe(({ policy, datasets, groups }) => {
        this.policy = policy;
        this.datasets = datasets;
        this.groups = groups;
        this.isLoading = false;
        this.cdr.markForCheck();
      });
  }

  duplicatePolicy() {
    this.policiesService
      .duplicateAgentPolicy(this.policyId || this.policy.id)
      .subscribe((resp) => {
        if (resp?.id) {
          this.notifications.success(
            'Agent Policy Duplicated',
            `New Agent Policy Name: ${resp?.name}`,
          );
          this.router.navigate([`view/${resp.id}`], {
            relativeTo: this.route.parent,
          });
        }
      });
  }

  ngOnDestroy() {
    this.policySubscription?.unsubscribe();
  }
}
