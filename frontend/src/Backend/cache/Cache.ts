class CacheEntry {
    value: any
    deadline: number

    constructor(value: any, ttl: number) {
        this.value = value
        this.deadline = Date.now() + 15 * 60 * 1000 // 15 minutes
    }
}

export default class Cache {
    data_ = new Map<string, CacheEntry>()

    constructor() {
        setInterval(() => this.clean(), 60 * 1000)
    }

    get<T>(key: string): T|null {
        const entry = this.data_.get(key)
        if (!entry) {
            return null
        }
        if (entry.deadline < Date.now()) {
            this.data_.delete(key)
            return null
        }
        return entry.value as T
    }

    set<T>(key: string, value: T, ttl: number = 15 * 60 * 1000) {
        this.data_.set(key, new CacheEntry(value, ttl))
    }

    delete(key: string) {
        this.data_.delete(key)
    }

    clean() {
        const now = Date.now()
        this.data_.forEach((v, k) => {
            if (v.deadline < now) {
                this.data_.delete(k)
            }
        })
    }

    clear() {
        this.data_.clear()
    }
}