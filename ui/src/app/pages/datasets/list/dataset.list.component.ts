import {
  AfterViewChecked,
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  OnInit,
  TemplateRef,
  ViewChild,
} from '@angular/core';

import { DropdownFilterItem } from 'app/common/interfaces/mainflux.interface';
import { ColumnMode, DatatableComponent, TableColumn } from '@swimlane/ngx-datatable';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { Debounce } from 'app/shared/decorators/utils';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { ActivatedRoute, Router } from '@angular/router';
import { DatasetDeleteComponent } from 'app/pages/datasets/delete/dataset.delete.component';
import { NbDialogService } from '@nebular/theme';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { DatasetDetailsComponent } from 'app/pages/datasets/details/dataset.details.component';

@Component({
  selector: 'ngx-dataset-list-component',
  templateUrl: './dataset.list.component.html',
  styleUrls: ['./dataset.list.component.scss'],
})
export class DatasetListComponent implements OnInit, AfterViewInit, AfterViewChecked {
  columnMode = ColumnMode;

  columns: TableColumn[];

  loading = false;

  paginationControls: OrbPagination<Dataset>;

  searchPlaceholder = 'Search by name';

  filterSelectedIndex = '0';

  // templates
  @ViewChild('actionsTemplateCell') actionsTemplateCell: TemplateRef<any>;

  tableFilters: DropdownFilterItem[] = [
    {
      id: '0',
      label: 'Name',
      prop: 'name',
      selected: false,
    },
  ];

  @ViewChild('tableWrapper') tableWrapper;

  @ViewChild(DatatableComponent) table: DatatableComponent;

  private currentComponentWidth;

  constructor(
    private cdr: ChangeDetectorRef,
    private dialogService: NbDialogService,
    private notificationsService: NotificationsService,
    private route: ActivatedRoute,
    private router: Router,
    private datasetPoliciesService: DatasetPoliciesService,
  ) {
    this.datasetPoliciesService.clean();
    this.paginationControls = DatasetPoliciesService.getDefaultPagination();
  }

  ngAfterViewChecked() {
    if (this.table && this.table.recalculate && (this.tableWrapper.nativeElement.clientWidth !== this.currentComponentWidth)) {
      this.currentComponentWidth = this.tableWrapper.nativeElement.clientWidth;
      this.table.recalculate();
      this.cdr.detectChanges();
      window.dispatchEvent(new Event('resize'));
    }
  }

  ngOnInit() {
    this.datasetPoliciesService.clean();
    this.getDatasets();
  }

  ngAfterViewInit() {
    this.columns = [
      {
        prop: 'name',
        name: 'Name',
        resizeable: false,
        flexGrow: 5,
        minWidth: 90,
      },
      {
        name: '',
        prop: 'actions',
        minWidth: 150,
        resizeable: false,
        sortable: false,
        flexGrow: 1,
        cellTemplate: this.actionsTemplateCell,
      },
    ];

    this.cdr.detectChanges();
  }


  @Debounce(500)
  getDatasets(pageInfo: NgxDatabalePageInfo = null): void {
    const isFilter = this.paginationControls.name?.length > 0 || this.paginationControls.tags?.length > 0;

    if (isFilter) {
      pageInfo = {
        offset: this.paginationControls.offset,
        limit: this.paginationControls.limit,
      };
      if (this.paginationControls.name?.length > 0) pageInfo.name = this.paginationControls.name;
      if (this.paginationControls.tags?.length > 0) pageInfo.tags = this.paginationControls.tags;
    }

    this.loading = true;
    this.datasetPoliciesService.getDatasetPolicies(pageInfo, isFilter).subscribe(
      (resp: OrbPagination<Dataset>) => {
        this.paginationControls = resp;
        this.paginationControls.offset = pageInfo?.offset || 0;
        this.paginationControls.total = resp.total;
        this.loading = false;
      },
    );
  }

  onOpenAdd() {
    this.router.navigate(['add'], {
      relativeTo: this.route.parent,
    });
  }

  onOpenEdit(dataset: any) {
    this.router.navigate(
      [`edit/${ dataset.id }`],
      {
        relativeTo: this.route.parent,
        state: { dataset: dataset, edit: true },
      },
    );
  }

  onFilterSelected(selectedIndex) {
    this.searchPlaceholder = `Search by ${ this.tableFilters[selectedIndex].label }`;
  }

  openDeleteModal(row: any) {
    const { id } = row;
    this.dialogService.open(DatasetDeleteComponent, {
      context: { name: row.name },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe(
      confirm => {
        if (confirm) {
          this.datasetPoliciesService.deleteDataset(id).subscribe(() => {
            this.getDatasets();
            this.notificationsService.success('Dataset successfully deleted', '');
          });
        }
      },
    );
  }

  openDetailsModal(row: any) {
    this.dialogService.open(DatasetDetailsComponent, {
      context: { dataset: row },
      autoFocus: true,
      closeOnEsc: true,
    }).onClose.subscribe((resp) => {
      if (resp) {
        this.onOpenEdit(row);
      } else {
        this.getDatasets();
      }
    });
  }

  searchDatasetItemByName(input) {
    this.getDatasets({
      ...this.paginationControls,
      [this.tableFilters[this.filterSelectedIndex].prop]: input,
    });
  }
}
