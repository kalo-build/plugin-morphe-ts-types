import { Person } from "./person"

export type ContactInfo = {
	email: string
	id: number
	personID?: number
	person?: Person
}

export type ContactInfoIDEmail = {
	email: string
}

export type ContactInfoIDPrimary = {
	id: number
}
