import { Person } from "./person"

export type Company = {
	id: number
	name: string
	taxID: string
	personIDs?: number[]
	persons?: Person[]
}

export type CompanyIDName = {
	name: string
}

export type CompanyIDPrimary = {
	id: number
}
