import { Variant, SortingOrder } from '@/room';

const sorters = new Map<SortingOrder, (a: Variant, b: Variant) => number>();

sorters.set(SortingOrder.DATE, (a: Variant, b: Variant) => b.createdAt - a.createdAt);

sorters.set(SortingOrder.RATING, (a: Variant, b: Variant) => {
  const difference = b.rating - a.rating;
  return difference !== 0 ? difference : b.createdAt - a.createdAt;
});

export default sorters;
