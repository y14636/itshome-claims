import { HomeComponent } from './home/home.component';
import { RouterModule, Routes } from "@angular/router";
import { NgModule } from '@angular/core';
import { AuthGuardService } from './auth-guard.service';
import { CallbackComponent } from './callback/callback.component';
import { ClaimsComponent } from './claims/claims.component';

const routes: Routes = [
  { path: '', redirectTo: 'home', pathMatch: 'full' },
  { path: 'home', component: HomeComponent },
  { path: 'claims', component: ClaimsComponent,  canActivate: [AuthGuardService] },
  { path: 'callback', component: CallbackComponent }
];

@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRoutingModule { }