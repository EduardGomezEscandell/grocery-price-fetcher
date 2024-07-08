export class AuthRefreshEndpoint {
    protected static path = '/api/auth/refresh'

    Path(): string {
        return AuthRefreshEndpoint.path
    }

    async POST(auth: string): Promise<string> {
        return fetch(AuthRefreshEndpoint.path, {
            method: 'POST',
            headers: {
                'Authorization': auth,
                'Accept': 'text/plain',
            },
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.text())
    }
}

export class MockAuthRefreshEndpoint extends AuthRefreshEndpoint {
    async POST(code: string): Promise<string> {
        console.log(`POST to ${this.Path()}: ${code}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => 'Mock Google endpoint')
    }
}