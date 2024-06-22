export default class ComparableString {
    displayName: string
    compareName: string

    constructor(displayName: string) {
        this.displayName = displayName
        this.compareName = this.localeFold(displayName)
    }

    private localeFold(s: string): string {
        return s.normalize("NFKD")           // Decompose unicode characters
            .replace(/[\u0300-\u036f]/g, "") // Remove accents
            .toLowerCase()
    }

    contains(other: ComparableString): boolean {
        return this.compareName.includes(other.compareName)
    }

    localeCompare(other: ComparableString): number {
        return this.compareName.localeCompare(other.compareName)
    }
}
