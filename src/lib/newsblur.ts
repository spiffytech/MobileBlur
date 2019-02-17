import axios from 'axios';

const base = 'https://newsblur.com';

export async function logIn(username: string, password: string) {
    const data = new FormData();
    data.set('username', username);
    data.set('password', password);
    const response = await axios({
        method: 'post',
        url: base + '/api/login',
        data,
    });
    return response.data;
}

export async function fetchFeeds() {
    const response = await axios.get(base + '/reader/feeds', {withCredentials: true});
    return response.data;
}

/**
 * Returns the stories for the given feed
 * @param feed Feed ID
 */
export async function fetchStories(feed: string | number, page: number) {
    const response = await axios.get(base + `/reader/feed/${feed}?page=${page}`);
    return response.data;
}

/**
 * Does this feed have any unread stories that match/exceed our current filter?
 * @param filter
 */
export function canShowFeed(filter: 'ng' | 'nt' | 'ps', {ng, nt, ps}: {ng: number, nt: number, ps: number}) {
    if (filter === 'ng') return ng > 0 || nt > 0 || ps > 0;
    if (filter === 'nt') return nt > 0 || ps > 0;
    if (filter === 'ps') return ps > 0;
    throw new Error(`Invalid filter level: '${filter}`);
}
