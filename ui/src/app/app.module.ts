/**
 * @license
 * Copyright Akveo. All Rights Reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 */
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { CoreModule } from './@core/core.module';
import { ThemeModule } from './@theme/theme.module';
import { AppComponent } from './app.component';
import { AppRoutingModule } from './app-routing.module';
import {
  NbAlertModule,
  NbButtonModule,
  NbCardModule,
  NbChatModule,
  NbCheckboxModule,
  NbDatepickerModule,
  NbDialogModule,
  NbIconModule,
  NbInputModule,
  NbLayoutModule,
  NbMenuModule,
  NbSidebarModule,
  NbToastrModule,
  NbWindowModule,
} from '@nebular/theme';

// MFx- Foorm dependency
import { FormsModule } from '@angular/forms';
// Mfx - MQTT dependencies for Gateways page
import { IMqttServiceOptions, MqttModule, MqttService } from 'ngx-mqtt';
import { environment } from 'environments/environment';
// Mfx - Auth and Profile pages
import { BreadcrumbModule } from 'xng-breadcrumb';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
import { ProfileComponent } from 'app/pages/profile/profile.component';
import { GoogleAnalyticsService } from './common/services/analytics/google-service-analytics.service';
import { MonacoEditorModule } from 'ngx-monaco-editor';

export const MQTT_SERVICE_OPTIONS: IMqttServiceOptions = {
  connectOnCreate: false,
  url: environment.mqttWsUrl,
};

@NgModule({
  declarations: [
    AppComponent,
    ProfileComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,

    AppRoutingModule,

    ThemeModule.forRoot(),
    NbSidebarModule.forRoot(),
    NbMenuModule.forRoot(),
    NbDatepickerModule.forRoot(),
    NbDialogModule.forRoot(),
    NbWindowModule.forRoot(),
    NbToastrModule.forRoot(),
    NbChatModule.forRoot({
      messageGoogleMapKey: 'AIzaSyA_wNuCzia92MAmdLRzmqitRGvCF7wCZPY',
    }),
    // Mfx dependencies
    CoreModule.forRoot(),
    MqttModule.forRoot(MQTT_SERVICE_OPTIONS),

    // 3rd party
    FormsModule,
    NbInputModule,
    NbCardModule,
    NbIconModule,
    NbButtonModule,
    BreadcrumbModule,
    NgxDatatableModule,
    NbAlertModule,
    NbCheckboxModule,
    NbLayoutModule,
    NbAlertModule,
    NbCheckboxModule,
    MonacoEditorModule.forRoot(),
  ],
  bootstrap: [AppComponent],
  // Mfx dependencies
  providers: [
    MqttService,
    GoogleAnalyticsService,
  ],
})
export class AppModule {
}
