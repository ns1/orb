import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  Input,
  OnChanges,
  OnDestroy,
  OnInit,
  Output,
  SimpleChanges,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';
import {
  ColumnMode,
  DatatableComponent,
  TableColumn,
} from '@swimlane/ngx-datatable';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { DatasetFromComponent, DATASET_RESPONSE } from 'app/pages/datasets/dataset-from/dataset-from.component';
import { DatasetDeleteComponent } from 'app/pages/datasets/delete/dataset.delete.component';
import { AgentGroupDetailsComponent } from 'app/pages/fleet/groups/details/agent.group.details.component';
import { SinkDetailsComponent } from 'app/pages/sinks/details/sink.details.component';
import { Subscription } from 'rxjs';
import { AgentMatchComponent } from 'app/pages/fleet/agents/match/agent.match.component';
import { OrbService } from 'app/common/services/orb.service';

@Component({
  selector: 'ngx-policy-datasets',
  templateUrl: './policy-datasets.component.html',
  styleUrls: ['./policy-datasets.component.scss'],
})
export class PolicyDatasetsComponent
  implements OnInit, OnDestroy, AfterViewInit, OnChanges {
  @Input()
  datasets: Dataset[];

  @Input()
  policy: AgentPolicy;

  isLoading: boolean;

  subscription: Subscription;

  errors;

  columnMode = ColumnMode;

  columns: TableColumn[];

  tableSorts = [
    {
      prop: 'name',
      dir: 'asc',
    },
  ];

  // templates
  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  @ViewChild('groupTemplateCell') groupTemplateCell: TemplateRef<any>;

  @ViewChild('validTemplateCell') validTemplateCell: TemplateRef<any>;

  @ViewChild('sinksTemplateCell') sinksTemplateCell: TemplateRef<any>;

  @ViewChild(DatatableComponent) table: DatatableComponent;

  private currentComponentWidth;

  constructor(
    private dialogService: NbDialogService,
    private cdr: ChangeDetectorRef,
    protected router: Router,
    protected route: ActivatedRoute,
    protected datasetService: DatasetPoliciesService,
    private notificationsService: NotificationsService,
    private orb: OrbService,
    ) {
    this.datasets = [];
    this.errors = {};
  }

  ngOnInit(): void {}

  ngOnChanges(changes: SimpleChanges) {}

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'agent_group',
        name: 'Agent Group',
        resizeable: false,
        canAutoResize: true,
        flexGrow: 2,
        cellTemplate: this.groupTemplateCell,
      },
      {
        prop: 'valid',
        name: 'Valid',
        resizeable: false,
        canAutoResize: true,
        flexGrow: 0.5,
        minWidth: 70,
        cellTemplate: this.validTemplateCell,
      },
      {
        prop: 'sinks',
        name: 'Sinks',
        resizeable: false,
        canAutoResize: true,
        flexGrow: 3,
        cellTemplate: this.sinksTemplateCell,
      },
      {
        name: '',
        prop: 'actions',
        resizeable: false,
        sortable: false,
        canAutoResize: true,
        flexGrow: 1,
        maxWidth: 130,
        minWidth: 130,
        cellTemplate: this.actionsTemplateCell,
      },
    ];

    this.cdr.detectChanges();
  }

  getTableHeight() {
    const rowHeight = 50;
    const headerHeight = 50;
    return (this.datasets.length * rowHeight + 15) + headerHeight + 'px';
  }
  onCreateDataset() {
    this.dialogService
      .open(DatasetFromComponent, {
        autoFocus: true,
        closeOnEsc: true,
        context: {
          policy: this.policy,
        },
        hasScroll: false,
        hasBackdrop: true,
        closeOnBackdropClick: true,
      })
      .onClose.subscribe((resp) => {
        if (resp === DATASET_RESPONSE.CREATED) {
          this.orb.refreshNow();
        }
      });
  }

  onOpenEdit(dataset) {
    this.dialogService
      .open(DatasetFromComponent, {
        autoFocus: true,
        closeOnEsc: false,
        context: {
          dataset,
          policy: this.policy,
        },
        hasScroll: false,
        closeOnBackdropClick: true,
        hasBackdrop: true,
      })
      .onClose.subscribe((resp) => {
        if (resp !== DATASET_RESPONSE.CANCELED) {
          this.orb.refreshNow();
        }
      });
  }

  onOpenEditAgentGroup(agentGroup: any) {
    this.router.navigate([`/pages/fleet/groups/edit/${agentGroup.id}`], {
      state: { agentGroup: agentGroup, edit: true },
      relativeTo: this.route,
    });
  }

  onOpenSinkDetails(sink) {
    this.dialogService
      .open(SinkDetailsComponent, {
        autoFocus: false,
        closeOnEsc: true,
        context: { sink },
        hasScroll: false,
      })
      .onClose.subscribe((resp) => {
        if (resp) {
          this.onOpenEditSink(sink);
        }
      });
  }

  onOpenEditSink(sink: any) {
    this.router.navigate([`pages/sinks/edit/${sink.id}`], {
      relativeTo: this.route,
      state: { sink: sink, edit: true },
    });
  }

  openDeleteModal(row: any) {
    const { name, id } = row;
    this.dialogService
      .open(DatasetDeleteComponent, {
        context: { name },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.datasetService.deleteDataset(id).subscribe(() => {
            this.notificationsService.success(
              'Dataset successfully deleted',
              '',
            );
          });
          this.orb.refreshNow();
        }
      });
  }

  ngOnDestroy() {
    this.subscription?.unsubscribe();
  }

  showAgentGroupMatches(agentGroup) {
    this.dialogService.open(AgentMatchComponent, {
      context: { agentGroup: agentGroup, policy: this.policy },
      autoFocus: true,
      closeOnEsc: true,
    });
  }

}
