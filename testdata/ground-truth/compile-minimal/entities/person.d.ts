import { Nationality } from "../enums/nationality"
import { Company } from "./company"

export type Person = {
	email: string
	id: number
	lastName: string
	nationality: Nationality
	companyID?: number
	company?: Company
}

export type PersonIDPrimary = {
	id: number
}
