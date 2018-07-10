import pify from 'pify';
import _jsonp from 'jsonp';

const jsonp = pify(_jsonp);

type AutocompleteResponse = [string, [string, 0, [number]][], any];
export async function googleAutocomplete(text: string): Promise<string[]> {
  const url = `https://www.google.ru/complete/search?gs_ri=psy-ab&q=${encodeURIComponent(text)}`;
  const response: AutocompleteResponse = await jsonp(url);

  return response[1].map(v => v[0]);
}

export async function yahooImages(query: string): Promise<string[]> {
  query = encodeURIComponent(query);

  const url = `https://cors-anywhere.herokuapp.com/https://images.search.yahoo.com/search/images?ei=UTF-8&p=${query}`;
  const pageText = await (await fetch(url)).text();
  const page = new DOMParser().parseFromString(pageText, 'text/html');

  const imageElements = Array.prototype.slice.call(page.querySelectorAll('.ld')) as HTMLLIElement[];

  const sources = imageElements
    .filter(e => e.id.startsWith('resitem-') && !e.id.startsWith('resitem-fpub'))
    .map(e => e.childNodes[0].childNodes[0] as HTMLImageElement)
    .map(e => e.src || e.dataset.src)
    .filter(e => e !== undefined) as string[];

  if (sources.length === 0) throw new Error('Unexpected yahoo images page layout');

  return sources;
}
