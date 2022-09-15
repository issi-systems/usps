export class AddressChecker {
    old: Map<string, Address>
    replace = new Map<string, Address>()

    constructor(public url: string, public src: Map<string, any>) {
        this.old = extractAddress(src)
    }
    get addresses(): Map<string, Address> {
        return extractAddress(this.src)
    }
    get changes(): Map<string, Address> {
        return diffAddress(this.old, this.addresses)
    }
    // call the go program to get a map of suggestions, check that these actually change.
    async suggest(address: Map<string, Address>): Promise<string[][]> {
        let js = JSON.stringify(Array.from(address.values()))
        console.log("suggest", js)

        if (address.size == 0) return [];
        this.replace = diffAddress(address, await rpc(this.url, js))
        let r: string[][] = []
        for (let [k, v] of this.replace) {
            r.push([addressToString(address.get(k)), addressToString(v)])
        }
        return r
    }

    // updates the local map and returns a map of changes for the server
    accept(): Map<string, string> {
        let r = new Map<string, string>()
        for (let [k, v] of this.replace) {
            let pc = (x: string, vs: string): void => {
                this.src.set(k + "-" + x, vs)
                r.set(k + "-" + x, vs)
            }
            pc("ADDRESS1", v.Address1)
            pc("ADDRESS2", v.Address2)
            pc("CITY", v.City)
            pc("STATE", v.State)
            pc("ZIP5", v.Zip5)
            pc("ZIP4", v.Zip4)
        }
        return r
    }
}

export interface Address {
    ID: string
    Address1: string
    Address2: string
    City: string
    State: string
    Zip5: string
    Zip4: string
    DPVConfirmation?: string
}

function addressToString(a: Address | null | undefined): string {
    if (!a) { return "" }
    return a.Address1
        + (a.Address2 ? "\n" + a.Address2 : "")
        + "\n" + a.City
        + ", " + a.State
        + " " + a.Zip5
        + (a.Zip4 ? "-" + a.Zip4 : "")

}
function compareAddress(a: Address, b: Address): number {
    return addressToString(a).localeCompare(addressToString(b))
}

// pull cobol addresses from string map into address map.
function extractAddress(m: Map<string, any>): Map<string, Address> {

    let r = new Map<string, Address>()
    for (let [k, v] of m) {
        if (k.indexOf("-CITY")) {
            let prefix = k.split("-")[0]

            let pc = (x: string): string => m.get(prefix + "-" + x) || ""
            r.set(prefix, {
                ID: prefix,
                Address1: pc("ADDRESS1"),
                Address2: pc("ADDRESS2"),
                City: pc("CITY"),
                State: pc("STATE"),
                Zip5: pc("ZIP"),
                Zip4: pc("ZIPEXT")
            })
        }
    }
    return r
}
// user can ignore, or set the address back into the map
// if they do, then we need to also send update methods to the core.


function diffAddress(old: Map<string, Address>, now: Map<string, Address>): Map<string, Address> {
    let r = new Map<string, Address>()
    for (let [k, v] of now) {
        let o = old.get(k)
        if (!o || 0 != compareAddress(v, o)) {
            r.set(k, v)
        }
    }
    return r
}


export function sm(x: Map<string, any>): string {
    return JSON.stringify(Object.fromEntries(x), null, 2)
}

async function rpc(url: string, body: string): Promise<Map<string, Address>> {
    let opt: RequestInit = {
        method: "POST",
        mode: 'cors',
        body: body,
        headers: {
            'content-type': 'application/json;charset=UTF-8',
        },
    }
    console.log("send", url, opt)
    let j = await fetch(url, opt);
    console.log(j)
    let js = await j.json() as Address[]
    console.log(js)
    let r = new Map<string, Address>()
    for (let o of js) {
        r.set(o.ID, o)
    }
    return r
}
