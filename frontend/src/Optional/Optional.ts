class Optional<T> {
    constructor(v: T|undefined) {
        this.value = v
    }

    private value: T|undefined;

    get(): T {
        if(this.value === undefined) {
            throw Error("Option is empty")
        }
        return this.value
    }

    hasValue(): boolean {
        return this.value !== undefined
    }

    then<U>(f: (t: T) => U|undefined): Optional<U> {
        if(this.value === undefined) {
            return new Optional<U>(undefined)
        }
        return new Optional<U>(f(this.value))
    }

    else(t: T): T {
        if(this.value === undefined) {
            return t
        }
        return this.value
    }
}

export default Optional