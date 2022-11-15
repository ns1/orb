import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';

@Component({
  selector: 'ngx-sink-details-component',
  templateUrl: './sink.details.component.html',
  styleUrls: ['./sink.details.component.scss'],
})
export class SinkDetailsComponent {
  strings = STRINGS.sink;

  @Input() sink: Sink = {};

  constructor(
    protected dialogRef: NbDialogRef<SinkDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
    !this.sink.tags ? this.sink.tags = {} : null;
  }

  onOpenEdit(sink: any) {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
