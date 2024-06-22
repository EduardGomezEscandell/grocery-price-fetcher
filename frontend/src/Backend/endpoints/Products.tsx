import { Product } from "../../State/State";
import Cache from "../cache/Cache";

export default class ProductsEndpoint {
    private path: string;
    protected cache: Cache | null = null;

    constructor(namespace: string, cache?: Cache) {
        this.path = `/api/products/${namespace}/`
        this.cache = cache || null;
    }

    Path(name?: string): string {
        return name === undefined
            ? this.path + '*'
            : this.path + name
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

    protected async get_one_uncached(name: string): Promise<Product> {
        return fetch(this.Path(name), {
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

    protected async post_uncached(oldName: string, p: Product): Promise<void> {
        return fetch(this.Path(oldName), {
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

    protected async delete_uncached(name: string): Promise<void> {
        return fetch(this.Path(name), {
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
        const cached = this.cache?.get<Product[]>(this.Path())
        if (cached) return cached

        return this
            .get_uncached()
            .then((products) => {
                this.cache?.set(this.Path(), products)
                return products
            })
    }

    async GET_ONE(name: string): Promise<Product> {
        const key = this.Path(name)
        const cached = this.cache?.get<Product>(key)
        if (cached) return cached

        return this
            .get_one_uncached(name)
            .then((product) => {
                this.cache?.set(key, product)
                return product
            })
    }

    async POST(oldName: string, product: Product): Promise<void> {
        this.cache?.delete(this.Path(oldName))
        return this
            .post_uncached(oldName, product)
            .then(() => { })
    }

    async DELETE(name: string): Promise<void> {
        this.cache?.delete(this.Path(name))
        return this.delete_uncached(name)
    }
}

export class MockProductsEndpoint extends ProductsEndpoint {
    constructor(namespace: string, cache?: Cache) {
        super(namespace, cache)
    }

    private static mockData = [
        { name: "Macarrons", price: 1.33, batch_size: 1, provider: 'Bonpreu', product_id: ['123', 'blabla'] },
        { name: "Ceba", price: 0.76, batch_size: 1, provider: 'Bonpreu', product_id: ['123', 'blabla'] },
        { name: "All", price: 0.88, batch_size: 3, provider: 'Mercadona', product_id: ['123', 'blabla'] },
        { name: "Tomàquet", price: 0.44, batch_size: 1, provider: 'Bonpreu', product_id: ['123', 'blabla'] },
        { name: "Oli", price: 0.2, batch_size: 1, provider: 'Bonpreu', product_id: ['123', 'blabla'] },
        { name: "Sal", price: 2.1, batch_size: 1, provider: 'Bonpreu', product_id: ['123', 'blabla'] },
        { name: "Pebre", price: 1.57, batch_size: 1, provider: 'Carrefour', product_id: ['123', 'blabla'] },
    ]

    protected async get_uncached(): Promise<Product[]> {
        console.log(`GET to ${this.Path()}`)
        return new Promise(resolve => setTimeout(resolve, 1000))
            .then(() => MockProductsEndpoint.mockData)
            .then((data: any[]) => data.map(Product.fromJSON))
    }

    protected async get_one_uncached(name: string): Promise<Product> {
        console.log(`GET to ${this.Path(name)}`)
        return new Promise(resolve => setTimeout(resolve, 1000))
            .then(() => MockProductsEndpoint.mockData)
            .then(Product.fromJSON)
    }

    protected async post_uncached(oldName: string, product: Product): Promise<void> {
        console.log(`POST to ${this.Path(oldName)}`)
        return new Promise(resolve => setTimeout(resolve, 1000))
            .then(() => { })
    }

    protected async delete_uncached(name: string): Promise<void> {
        if (name === 'Sal') {
            return new Promise(resolve => setTimeout(resolve, 1000))
                .then(() => Promise.reject(new Response('Bla bla bla terrible error', { status: 500 })))
        }

        console.log(`DELETE to ${this.Path(name)}`)
        return new Promise(resolve => setTimeout(resolve, 1000))
            .then(() => { })
    }
}