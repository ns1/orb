import { NgModule } from '@angular/core';
import {
  NbButtonModule,
  NbCardModule,
  NbCheckboxModule,
  NbDialogService,
  NbInputModule,
  NbListModule,
  NbMenuModule,
  NbSelectModule,
  NbTabsetModule,
  NbWindowService,
} from '@nebular/theme';
import { ThemeModule } from '../@theme/theme.module';
import { PagesComponent } from './pages.component';
import { DashboardModule } from './dashboard/dashboard.module';
import { PagesRoutingModule } from './pages-routing.module';

// Mainflux - Dependencies
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
// Mainflux - Common and Shared
import { SharedModule } from 'app/shared/shared.module';
import { CommonModule } from 'app/common/common.module';
import { ConfirmationComponent } from 'app/shared/components/confirmation/confirmation.component';

// ORB
import { DatasetsComponent } from 'app/pages/datasets/datasets.component';
import { DatasetsAddComponent } from 'app/pages/datasets/add/datasets.add.component';
import { DatasetsDetailsComponent } from 'app/pages/datasets/details/datasets.details.component';
import { FleetsComponent } from 'app/pages/fleets/fleets.component';
import { FleetsAddComponent } from 'app/pages/fleets/add/fleets.add.component';
import { FleetsDetailsComponent } from 'app/pages/fleets/details/fleets.details.component';
import { SinksComponent } from 'app/pages/sinks/sinks.component';
import { SinksAddComponent } from 'app/pages/sinks/add/sinks.add.component';
import { SinksDetailsComponent } from 'app/pages/sinks/details/sinks.details.component';
import { SinksDeleteComponent } from 'app/pages/sinks/delete/sinks.delete.component';
import { AgentsComponent } from 'app/pages/agents/agents.component';
import { AgentsAddComponent } from 'app/pages/agents/add/agents.add.component';
import { AgentsDetailsComponent } from 'app/pages/agents/details/agents.details.component';
import { MatInputModule } from '@angular/material/input';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';

@NgModule({
  imports: [
    PagesRoutingModule,
    ThemeModule,
    NbMenuModule,
    DashboardModule,
    SharedModule,
    CommonModule,
    FormsModule,
    NbButtonModule,
    NbCardModule,
    NbInputModule,
    NbSelectModule,
    NbCheckboxModule,
    NbListModule,
    NbTabsetModule,
    ReactiveFormsModule,
    MatInputModule,
    MatChipsModule,
    MatIconModule,
  ],
  exports: [
    SharedModule,
    CommonModule,
    FormsModule,
    NbButtonModule,
    NbCardModule,
    NbInputModule,
    NbSelectModule,
    NbCheckboxModule,
    NbListModule,
  ],
  declarations: [
    PagesComponent,
    // Orb
    // Agent Group Management
    AgentsComponent,
    AgentsAddComponent,
    AgentsDetailsComponent,
    // Dataset Explorer
    DatasetsComponent,
    DatasetsAddComponent,
    DatasetsDetailsComponent,
    // Fleet Management
    FleetsComponent,
    FleetsAddComponent,
    FleetsDetailsComponent,
    // Sink Management
    SinksComponent,
    SinksAddComponent,
    SinksDetailsComponent,
    SinksDeleteComponent,
  ],
  providers: [
    NbDialogService,
    NbWindowService,
  ],
  entryComponents: [
    ConfirmationComponent,
  ],
})
export class PagesModule {
}
