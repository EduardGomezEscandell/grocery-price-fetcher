import { Product } from "../../State/State";
import Cache from "../cache/Cache";

export default class ProviderEndpoint {
    private auth: string
    cache: Cache | null = null;
    private static path = `/api/provider`

    constructor(auth: string, cache?: Cache) {
        this.auth = auth
        this.cache = cache || null;
    }

    static Path(q: Query): string {
        return ProviderEndpoint.path + encodeURI(`?provider=${q.provider}&id=${q.product_code}`)
    }

    protected async get_uncached(q: Query): Promise<number> {
        return fetch(ProviderEndpoint.Path(q), {
            method: 'GET',
            headers: {
                'Authorization': this.auth,
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.json())
            .then(v => v.price)
    }

    async GET(p: Product): Promise<number> {
        const key = ProviderEndpoint.Path(p)
        const cached = this.cache?.get<number>(key)
        if (cached) return cached

        return this
            .get_uncached(p)
            .then((resp) => {
                this.cache?.set(key, resp)
                return resp
            })
    }
}

interface Query {
    provider: string;
    product_code: string;
}

export class MockProvidersEndpoint extends ProviderEndpoint {
    protected async get_uncached(q: Query): Promise<number> {
        const path = ProviderEndpoint.Path(q)
        console.log(`GET to ${path}`)

        if (q.product_code === '404') {
            // Mock a 404
            return new Promise(resolve => setTimeout(resolve, 1000))
                .then(() => Promise.reject(new Response('Not Found', { status: 404 })))
        }

        return new Promise(resolve => setTimeout(resolve, 1000))
            .then(() => {return 11.116}) // Extra precision to test rounding
    }
}