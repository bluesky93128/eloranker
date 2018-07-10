import { EventEmitter } from 'events';
import * as events from './events';
import { Variant, EditMode } from './room';

function debug(...args: any[]) {
  // if (process.env.NODE_ENV !== 'production') {
  console.log(...args);
  // }
}

export interface Connection extends EventEmitter {
  on(event: 'state', listener: (state: number) => void): this;
  once(event: 'state', listener: (state: number) => void): this;
  on(event: 'error', listener: (error?: string) => void): this;
  once(event: 'error', listener: (error?: string) => void): this;

  on(event: 'room:new', listener: (event: events.NewRoomEvent) => void): this;
  once(event: 'room:new', listener: (event: events.NewRoomEvent) => void): this;
  on(event: 'room:join', listener: (event: events.JoinRoomEvent) => void): this;
  once(event: 'room:join', listener: (event: events.JoinRoomEvent) => void): this;
  on(event: 'variant:allocate', listener: (event: events.AllocateVariantEvent) => void): this;
  once(event: 'variant:allocate', listener: (event: events.AllocateVariantEvent) => void): this;
  on(event: 'voting:get', listener: (event: events.GetVotingEvent) => void): this;
  once(event: 'voting:get', listener: (event: events.GetVotingEvent) => void): this;

  on(event: 'room:clients', listener: (event: events.RoomClientsEvent) => void): this;
  once(event: 'room:clients', listener: (event: events.RoomClientsEvent) => void): this;
  on(event: 'variant:update', listener: (event: events.UpdateVariantEvent) => void): this;
  once(event: 'variant:update', listener: (event: events.UpdateVariantEvent) => void): this;
  on(event: 'variant:remove', listener: (event: events.RemoveVariantEvent) => void): this;
  once(event: 'variant:remove', listener: (event: events.RemoveVariantEvent) => void): this;

  on(event: 'settings:title', listener: (event: events.SettingsTitleEvent) => void): this;
  once(event: 'settings:title', listener: (event: events.SettingsTitleEvent) => void): this;
  on(event: 'settings:quotaEnabled', listener: (event: events.SettingsQuotaEvent) => void): this;
  once(event: 'settings:quotaEnabled', listener: (event: events.SettingsQuotaEvent) => void): this;
  on(event: 'settings:editMode', listener: (event: events.SettingsEditModeEvent) => void): this;
  once(event: 'settings:editMode', listener: (event: events.SettingsEditModeEvent) => void): this;
}

export class Connection extends EventEmitter {
  private ws: WebSocket;
  constructor(url: string) {
    super();
    this.ws = new WebSocket(url);

    this.ws.onopen = () => this.emit('state', this.ws.readyState);
    this.ws.onclose = () => this.emit('state', this.ws.readyState);
    this.ws.onerror = () => this.emit('error');

    this.ws.onmessage = ev => {
      (ev.data as string).split('\n').forEach(data => {
        const response = JSON.parse(data);
        debug('[SOCKET:GET]', response);
        if (!response.event) {
          this.emit('error', response.error != null ? response.error : 'Invalid response');
          return;
        }
        if (this.listenerCount(response.event) === 0) {
          debug('WARNING: No listeners for this event ^');
          return;
        }

        this.emit(response.event, response);
      });
    };
  }

  private send(type: string, data?: { [key: string]: any }) {
    debug('[SOCKET:SND]', { type, data });
    this.ws.send(JSON.stringify({ type, data }));
  }

  public async waitOpen() {
    if (this.ws.readyState === WebSocket.OPEN) return Promise.resolve();
    if (this.ws.readyState >= WebSocket.CLOSING) throw new Error('Socket closed');

    return new Promise(resolve => {
      const handler = (state: number) => {
        if (state === WebSocket.OPEN) {
          resolve();
          this.removeListener('state', handler);
        } else if (state >= WebSocket.CLOSING) {
          throw new Error('Socket closed');
        }
      };
      this.on('state', handler);
    });
  }

  public newRoom(title: string) {
    this.send('room:new', { title });
    return new Promise<events.NewRoomEvent>((resolve, reject) =>
      this.once(
        'room:new',
        event => (event.error ? reject(new Error(event.error)) : resolve(event)),
      ),
    );
  }

  public joinRoom(name: string, secret: string) {
    this.send('room:join', { name, secret });
    return new Promise<events.JoinRoomEvent>((resolve, reject) =>
      this.once(
        'room:join',
        event => (event.error ? reject(new Error(event.error)) : resolve(event)),
      ),
    );
  }

  public leaveRoom() {
    this.send('room:leave');
  }

  public allocateNewVariant() {
    this.send('variant:allocate');
    return new Promise<string>((resolve, reject) =>
      this.once(
        'variant:allocate',
        event => (event.error ? reject(new Error(event.error)) : resolve(event.uuid)),
      ),
    );
  }

  public updateVariant(variant: Variant) {
    this.send('variant:update', { uuid: variant.uuid, text: variant.text, image: variant.image });
  }

  public getVoting() {
    this.send('voting:get');
  }

  public submitVote(uuid: string) {
    this.send('voting:submit', { uuid });
  }

  public setVariantIgnored(uuid: string, ignored: boolean) {
    this.send('variant:setIgnored', { uuid, ignored });
  }

  public removeVariant(id: string) {
    this.send('variant:remove', { id });
  }

  public importVariants(owner: string, repo: string) {
    this.send('variant:import', { type: 'github-issues', data: { owner, repo } });
  }

  public setTitle(value: string) {
    this.send('settings:title', { value });
  }

  public setQuotaEnabled(value: boolean) {
    this.send('settings:quotaEnabled', { value });
  }

  public setEditMode(value: EditMode) {
    this.send('settings:editMode', { value });
  }
}

const socketHost =
  process.env.VUE_APP_SOCKET_HOST ||
  (process.env.NODE_ENV === 'production' ? window.location.host : 'localhost');

const socketUrl = `${window.location.protocol === 'https' ? 'wss' : 'ws'}://${socketHost}/ws`;

export default new Connection(socketUrl);
