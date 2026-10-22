package ipam

// IPPool is a share of a subnet's addresses, not a separate source of them. Every Subnet
// holds its pools under Subnet.IPPools, where the entry named "" is the default pool and
// starts out owning the subnet's whole usable range; creating a named pool moves the
// addresses it consists of from the default pool to it, see Subnet.setIPPoolLocked.
//
// Which addresses a consumer can get depends on how it asks:
//   - a random allocation is scoped to one pool - Subnet.GetRandomAddress("") hands out
//     default pool addresses only, and a named pool hands out only its own. This is what
//     a pool is for: it narrows "whoever uses this subnet" down to "whoever references
//     this pool";
//   - a static allocation is not scoped at all - Subnet.GetStaticAddress looks the
//     address up in every pool and allocates from whichever one owns it, so a consumer
//     naming an address explicitly can take it out of any pool.
//
// Three fields of this struct sound alike and mean different things:
//   - V4IPs/V6IPs are the addresses the pool consists of;
//   - V4Reserved/V6Reserved are the ones among them the subnet excludes (excludeIps).
//     Excluded addresses are in no free list, so nobody is handed one at random, but a
//     static allocation still takes them - see the excludeIps branch of GetStaticAddress;
//   - PendingActivation marks a membership that is carved out of the subnet but not
//     handed out yet, the state a pool is in between being reserved and being activated,
//     see Subnet.ReserveIPPoolIPs. It stops random allocations only: a static allocation
//     keeps resolving, which is how the startup restore records the addresses the pool's
//     pods already hold.
type IPPool struct {
	PendingActivation bool
	V4IPs             *IPRangeList
	V4Free            *IPRangeList
	V4Available       *IPRangeList
	V4Reserved        *IPRangeList
	V4Released        *IPRangeList
	V4Using           *IPRangeList
	V6IPs             *IPRangeList
	V6Free            *IPRangeList
	V6Available       *IPRangeList
	V6Reserved        *IPRangeList
	V6Released        *IPRangeList
	V6Using           *IPRangeList
}
