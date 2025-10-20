import { Company } from "./company"
import { Person } from "./person"

export type Comment = {
	id: number
	text: string
	commentableID?: string
	commentableType?: string
	commentable?: Person | Company
}

export type CommentIDPrimary = {
	id: number
}
