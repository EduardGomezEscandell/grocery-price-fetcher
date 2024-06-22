import { Product } from "../../State/State";
import Cache from "../cache/Cache";

export default class ProductsEndpoint {
    private path: string;
    protected cache: Cache | null = null;

    constructor(namespace: string, cache?: Cache) {
        this.path = `/api/products/${namespace}/`
        this.cache = cache || null;
    }

    PathAll(): string {
        return this.path + '*'
    }

    Path(id: number): string {
        return this.path + id.toString()
    }

    protected async get_uncached(): Promise<Product[]> {
        return fetch(this.PathAll(), {
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

    protected async get_one_uncached(id: number): Promise<Product> {
        return fetch(this.Path(id), {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.json())
            .then(Product.fromJSON)
    }

    protected async post_uncached(p: Product): Promise<void> {
        return fetch(this.Path(p.id), {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
            body: JSON.stringify({...p, product_id: [p.product_id, "", ""]}),
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(() => { })
    }

    protected async delete_uncached(id: number): Promise<void> {
        return fetch(this.Path(id), {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            },
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(() => { })
    }

    async GET(): Promise<Product[]> {
        const cached = this.cache?.get<Product[]>(this.PathAll())
        if (cached) return cached

        return this
            .get_uncached()
            .then((products) => {
                this.cache?.set(this.PathAll(), products)
                return products
            })
    }

    async GET_ONE(id: number): Promise<Product> {
        const key = this.Path(id)
        const cached = this.cache?.get<Product>(key)
        if (cached) return cached

        return this
            .get_one_uncached(id)
            .then((product) => {
                this.cache?.set(key, product)
                return product
            })
    }

    async POST(product: Product): Promise<void> {
        this.cache?.delete(this.Path(product.id))
        return this
            .post_uncached(product)
            .then(() => { })
    }

    async DELETE(id: number): Promise<void> {
        this.cache?.delete(this.Path(id))
        return this.delete_uncached(id)
    }
}

export class MockProductsEndpoint extends ProductsEndpoint {
    constructor(namespace: string, cache?: Cache) {
        super(namespace, cache)
    }

    private static mockData = [
        { id: 1, name: "Macarrons", price: 1.33, batch_size: 1, provider: 'Bonpreu', product_id: ['123', 'blabla'] },
        { id: 2, name: "Ceba", price: 0.76, batch_size: 1, provider: 'Bonpreu', product_id: ['123', 'blabla'] },
        { id: 3, name: "All", price: 0.88, batch_size: 3, provider: 'Mercadona', product_id: ['123', 'blabla'] },
        { id: 4, name: "Tomàquet", price: 0.44, batch_size: 1, provider: 'Bonpreu', product_id: ['123', 'blabla'] },
        { id: 5, name: "Oli", price: 0.2, batch_size: 1, provider: 'Bonpreu', product_id: ['123', 'blabla'] },
        { id: 404, name: "Sal", price: 2.1, batch_size: 1, provider: 'Bonpreu', product_id: ['123', 'blabla'] },
        { id: 6, name: "Pebre", price: 1.57, batch_size: 1, provider: 'Carrefour', product_id: ['123', 'blabla'] },
    ]

    protected async get_uncached(): Promise<Product[]> {
        console.log(`GET to ${this.PathAll()}`)
        return new Promise(resolve => setTimeout(resolve, 1000))
            .then(() => MockProductsEndpoint.mockData)
            .then((data: any[]) => data.map(Product.fromJSON))
    }

    protected async get_one_uncached(id: number): Promise<Product> {
        console.log(`GET to ${this.Path(id)}`)
        return new Promise(resolve => setTimeout(resolve, 1000))
            .then(() => MockProductsEndpoint.mockData)
            .then(Product.fromJSON)
    }

    protected async post_uncached(product: Product): Promise<void> {
        console.log(`POST to ${this.Path(product.id)}`)
        if (product.id === 404) {
            return new Promise(resolve => setTimeout(resolve, 1000))
                .then(() => Promise.reject(new Response('Bla bla bla terrible error', { status: 500 })))
        }

        return new Promise(resolve => setTimeout(resolve, 1000))
            .then(() => { })
    }

    protected async delete_uncached(id: number): Promise<void> {
        console.log(`DELETE to ${this.Path(id)}`)
        if (id === 404) {
            return new Promise(resolve => setTimeout(resolve, 1000))
                .then(() => Promise.reject(new Response('Bla bla bla terrible error', { status: 500 })))
        }

        return new Promise(resolve => setTimeout(resolve, 1000))
            .then(() => { })
    }
}