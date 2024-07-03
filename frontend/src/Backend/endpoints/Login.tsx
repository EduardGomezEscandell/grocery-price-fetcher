export class LoginEndpoint {
    protected static path = '/api/login'

    Path(): string {
        return LoginEndpoint.path
    }

    async POST(code: string): Promise<string> {
        return fetch(LoginEndpoint.path, {
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

export class MockLoginEndpoint extends LoginEndpoint {
    async POST(code: string): Promise<string> {
        console.log(`POST to ${this.Path()}: ${code}`)
        return new Promise(resolve => setTimeout(resolve, 100))
            .then(() => 'Mock Google endpoint')
    }
}