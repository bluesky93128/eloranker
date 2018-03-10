export enum EditMode {
  Trust,
  Normal,
  Restricted,
}

export interface Variant {
  uuid: string;
  text: string;
  image: string;
  rating: number;
  createdAt: number;
  author: string;
}

const variantDefaults: Variant = {
  uuid: '',
  text: '',
  image: '',
  rating: 0,
  createdAt: -1,
  author: '',
};
export function emptyVariant(init: Partial<Variant> = {}) {
  return Object.assign({}, variantDefaults, init);
}

export enum SortingOrder {
  DATE,
  RATING,
}
