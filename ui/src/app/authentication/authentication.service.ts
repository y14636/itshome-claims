import { Injectable } from '@angular/core';
import { ILogin } from '../login';
import { BehaviorSubject } from 'rxjs';
import { Router } from '@angular/router';

// const clientSecret = "K4pH3cT7fR1lX6pM2uN1tC4iA6rJ7jA7rS0sO6mK3vU0sI7iL0"; //dev
// const clientId = "6d888891-0f0a-470e-a351-95d422c66513"; // dev

@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {
  private loggedIn = new BehaviorSubject<boolean>(false);

  get isLoggedIn() {
    return this.loggedIn.asObservable(); // {2}
  }

  constructor(private router: Router) { }

  login(user: ILogin){
    if (user.userid !== '' && user.password !== '' ) { // {3}
    console.log("inside login");
      this.loggedIn.next(true);
      this.router.navigate(['/']);
    }
  }

  logout(): void {
    this.loggedIn.next(false);
    localStorage.setItem('isLoggedIn', "false");
    localStorage.removeItem('token');
    this.router.navigate(['/login']);
  } 

}