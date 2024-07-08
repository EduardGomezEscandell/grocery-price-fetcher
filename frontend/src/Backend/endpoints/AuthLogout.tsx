export class AuthLogoutEndpoint {
    protected static path = '/api/auth/logout'
    private auth: string

    constructor(auth: string) {
        this.auth = auth
    }

    Path(): string {
        return AuthLogoutEndpoint.path
    }

    async POST(): Promise<void> {
        return fetch(AuthLogoutEndpoint.path, {
            method: 'POST',
            headers: {
                'Authorization': this.auth,
            },
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(() => { })
    }
}

export class MockAuthLogoutEndpoint extends AuthLogoutEndpoint {
    async POST(): Promise<void> {
        console.log(`POST to ${this.Path()}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => { })
    }
}