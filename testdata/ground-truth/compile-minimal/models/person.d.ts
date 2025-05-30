import { Nationality } from "../enums/nationality"
import { Company } from "./company"
import { ContactInfo } from "./contact-info"

export type Person = {
	firstName: string
	id: number
	lastName: string
	nationality: Nationality
	companyID?: number
	company?: Company
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
