import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { AuthenticationService } from './authentication/authentication.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit{
  //title = 'ui';
  isLoggedIn$: Observable<boolean>;

  constructor(private authenticationService: AuthenticationService) { }

  ngOnInit() {
    this.isLoggedIn$ = this.authenticationService.isLoggedIn;
  }

  onLogout(){
    this.authenticationService.logout();
  }
}