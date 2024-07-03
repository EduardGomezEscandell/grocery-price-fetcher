export class LogoutEndpoint {
    protected static path = '/api/logout'
    private auth: string

    constructor(auth: string) {
        this.auth = auth
    }

    Path(): string {
        return LogoutEndpoint.path
    }

    async POST(): Promise<void> {
        return fetch(LogoutEndpoint.path, {
            method: 'POST',
            headers: {
                'Authorization': this.auth,
                'Content-Type': 'application/json',
                'Accept': 'text/plain',
            },
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(() => { })
    }
}

export class MockLogoutEndpoint extends LogoutEndpoint {
    async POST(): Promise<void> {
        console.log(`POST to ${this.Path()}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => { })
    }
}