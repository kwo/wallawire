export interface IAuthResponse {
    status: AuthStatus;
    message: string;
}

export enum AuthStatus {
    Success = 1,
    BadRequest,
    Unauthorized,
    Forbidden,
    Unknown,
    Error,
}
