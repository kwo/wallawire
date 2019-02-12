export function to(p: Promise<any>): Promise<[any, Error | null]> {
    return p.then((data: any): [any, Error | null] => {
       return [data, null];
    }).catch((err: Error): [any, Error | null] => [null, err]);
 }
