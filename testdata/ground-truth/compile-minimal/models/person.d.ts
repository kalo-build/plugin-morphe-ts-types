import { Nationality } from "../enums/nationality"
import { Comment } from "./comment"
import { Company } from "./company"
import { Contact } from "./contact"
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
	noteIDs?: number[]
	notes?: Comment[]
	personalContactID?: number
	personalContact?: Contact
	workContactID?: number
	workContact?: Contact
}

export type PersonIDName = {
	firstName: string
	lastName: string
}

export type PersonIDPrimary = {
	id: number
}
