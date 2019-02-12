export interface IChangeProfileResponse {
    status: ChangeProfileStatus;
    message: string;
}

export enum ChangeProfileStatus {
    Success = 1,
    BadRequest,
    Unauthorized,
    Forbidden,
    Error,
    Unknown,
}
