import { Product } from "../../State/State";
import Cache from "../cache/Cache";

export default class ProductsEndpoint {
    path: string;
    cache: Cache | null = null;

    constructor(namespace: string, cache?: Cache) {
        this.path = `/api/products/${namespace}`
        this.cache = cache || null;
    }

    Path(): string {
        return this.path;
    }

    protected async get_uncached(): Promise<Product[]> {
        return fetch(this.Path(), {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.json())
            .then((data: any[]) => data.map(Product.fromJSON))
    }

    async GET(): Promise<Product[]> {
        const cached = this.cache?.get<Product[]>(this.path)
        if (cached) return cached

        return this
            .get_uncached()
            .then((products) => {
                this.cache?.set(this.Path(), products)
                return products
            })
    }
}

export class MockProductsEndpoint extends ProductsEndpoint {
    constructor(namespace: string, cache?: Cache) {
        super(namespace, cache)
    }

    protected async get_uncached(): Promise<Product[]> {
        console.log(`GET to ${this.path}:`)
        return new Promise(resolve => setTimeout(resolve, 1000))
            .then(() => [
                { name: "Macarrons", price: 1.33, batch_size: 1, provider: 'Bonpreu' },
                { name: "Ceba", price: 0.76, batch_size: 1, provider: 'Bonpreu' },
                { name: "All", price: 0.88, batch_size: 3, provider: 'Mercadona' },
                { name: "Tomàquet", price: 0.44, batch_size: 1, provider: 'Bonpreu'},
                { name: "Oli", price: 0.2, batch_size: 1, provider: 'Bonpreu' },
                { name: "Sal", price: 2.1, batch_size: 1, provider: 'Bonpreu' },
                { name: "Pebre", price: 1.57, batch_size: 1, provider: 'Carrefour' }
            ])
            .then((data: any[]) => data.map(Product.fromJSON))
    }
}