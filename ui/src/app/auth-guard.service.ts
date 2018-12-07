import { CanActivate, CanActivateChild, Router, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Injectable } from '@angular/core';
import { AuthenticationService } from './authentication/authentication.service';
import { Observable } from 'rxjs';
import { map, take } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class AuthGuardService implements CanActivate {

  constructor(private authService: AuthenticationService,private router: Router) { }
  
  // canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean {
  //   let url: string = state.url;  
  //   console.log("canActivate", this.verifyLogin(url));
  //   return this.verifyLogin(url);
  canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): Observable<boolean> {
    return this.authService.isLoggedIn         // {1}
      .pipe(
        take(1),                              // {2} 
        map((isLoggedIn: boolean) => {         // {3}
          if (!isLoggedIn){
            this.router.navigate(['/login']);  // {4}
            return false;
          }
          return true;
        })
      )
  }
}

//   verifyLogin(url) : boolean{
//       if(!this.isLoggedIn()){
//           this.router.navigate(['/login']);
//           return false;
//       }
//       else if(this.isLoggedIn()){
//           return true;
//       }
//   }
//   public isLoggedIn(): boolean{
//       let status = false;
//       if( localStorage.getItem('isLoggedIn') == "true"){
//         status = true;
//       }
//       else{
//         status = false;
//       }
//       return status;
//   }
// }