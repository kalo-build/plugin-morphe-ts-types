import { Nationality } from "../enums/nationality"
import { ContactInfo } from "./contact-info"

export type Person = {
	firstName: string
	id: number
	lastName: string
	nationality: Nationality
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
