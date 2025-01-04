import { ContactInfo } from "./contact-info"

export type Person = {
	firstName: string
	id: number
	lastName: string
	contactInfoID?: number
	contactInfo?: ContactInfo
}

export type PersonIDName = {
	firstName: string
	lastName: string
}

export type PersonIDPrimary = {
	id: number
}
