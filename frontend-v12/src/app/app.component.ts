import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'my-app';

  clickMessage = '';

  onClickMe() {
    // @ts-ignore
    window.backend.basic("TEST").then(result =>
      this.clickMessage = result
    );
  }
}
