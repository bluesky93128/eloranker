import { Variant, EditMode } from './room';

interface ResponseEvent {
  event: string;
  error?: string;
}

export type NewRoomEvent = ResponseEvent & {
  name: string;
  secret: string;
};

export type JoinRoomEvent = ResponseEvent & {
  variants: Variant[];
  identifier: string;

  isAdmin: boolean;
  title: string;
  quotaEnabled: boolean;
  editMode: EditMode;
  ignoredVariants: { [id: string]: true };
};

export type AllocateVariantEvent = ResponseEvent & Partial<Variant>;

export type GetVotingEvent = ResponseEvent & {
  variants: [string, string];
  // quota: { current: number; limit: number };
};

export type RoomClientsEvent = {
  event: string;
  clients: number;
};

export type UpdateVariantEvent = Partial<Variant> & {
  event: string;
  error?: string;
  uuid: string;
};

export type RemoveVariantEvent = {
  event: string;
  error?: string;
  id: string;
};

export type SettingsTitleEvent = ResponseEvent & { value: string };
export type SettingsQuotaEvent = ResponseEvent & { value: boolean };
export type SettingsEditModeEvent = ResponseEvent & { value: EditMode };
