export class AuthLoginEndpoint {
    protected static path = '/api/auth/login'

    Path(): string {
        return AuthLoginEndpoint.path
    }

    async POST(code: string): Promise<string> {
        return fetch(AuthLoginEndpoint.path, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'text/plain',
            },
            body: JSON.stringify({
                code: code
            })
        })
            .then(r => r.ok ? r : Promise.reject(r))
            .then(r => r.text())
    }
}

export class MockAuthLoginEndpoint extends AuthLoginEndpoint {
    async POST(code: string): Promise<string> {
        console.log(`POST to ${this.Path()}: ${code}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => 'Bearer MOCK_123456789')
    }
}