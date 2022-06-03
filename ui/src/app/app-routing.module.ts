import { ExtraOptions, RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { AuthGuard } from 'app/auth/auth-guard.service';

// Mfx- Custom Logout and Register components that
// replace NbLogoutComponent and NbRegisterComponent

export const routes: Routes = [
  {
    path: 'pages',
    loadChildren: () => import('./pages/pages.module')
      .then(m => m.PagesModule),
    data: {breadcrumb: {skip: true}},
    canActivate: [AuthGuard],
  },
  {
    path: 'auth',
    loadChildren: () => import('./auth/auth.module')
      .then(m => m.AuthModule),
    data: {breadcrumb: {skip: true}},
  },
  {path: '', redirectTo: 'pages/home', pathMatch: 'full'},
  {path: '**', redirectTo: 'pages/home'},
];

const config: ExtraOptions = {
  useHash: false,
};

@NgModule({
  imports: [RouterModule.forRoot(routes, config)],
  exports: [RouterModule],
})
export class AppRoutingModule {
}
