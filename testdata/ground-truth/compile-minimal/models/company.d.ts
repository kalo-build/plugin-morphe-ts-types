import { Comment } from "./comment"
import { Contact } from "./contact"
import { Person } from "./person"

export type Company = {
	id: number
	name: string
	taxID: string
	mailingContactID?: number
	mailingContact?: Contact
	mainContactID?: number
	mainContact?: Contact
	noteIDs?: number[]
	notes?: Comment[]
	personIDs?: number[]
	persons?: Person[]
}

export type CompanyIDName = {
	name: string
}

export type CompanyIDPrimary = {
	id: number
}
