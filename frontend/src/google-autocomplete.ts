import pify from 'pify';
import _jsonp from 'jsonp';

const jsonp = pify(_jsonp);

type AutocompleteResponse = [string, [string, 0, [number]][], any];
export default async function googleAutocomplete(text: string): Promise<string[]> {
  const url = `https://www.google.ru/complete/search?gs_ri=psy-ab&q=${encodeURIComponent(text)}`;
  const response: AutocompleteResponse = await jsonp(url);

  return response[1].map(v => v[0]);
}
