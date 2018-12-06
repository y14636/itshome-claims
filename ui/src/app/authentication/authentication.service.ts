import { Injectable } from '@angular/core';
import { ILogin } from '../login';

// const clientSecret = "K4pH3cT7fR1lX6pM2uN1tC4iA6rJ7jA7rS0sO6mK3vU0sI7iL0"; //dev
// const clientId = "6d888891-0f0a-470e-a351-95d422c66513"; // dev

@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {

  constructor() { }

  logout(): void {
    localStorage.setItem('isLoggedIn', "false");
    localStorage.removeItem('token');
  } 

}