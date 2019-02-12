export interface IServerStatus {
    service: string;
    version: string;
    runtime: string;
    buildTime: Date;
    start: Date;
    time: Date;
    uptime: string;
}
