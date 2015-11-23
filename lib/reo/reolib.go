package reo

// FIXME maybe we need to send a stop signal to any potential
// blocking operation? any SyncRead may leads to this kind of bugs
// NOTE when later some more complicated example would be trapped
// into deadlocks plz kindly check this
// FIXME maybe there're lots of WaitRead that need to replaced by
// select ...

func SyncChannel(in, out, stop Port) {
	for {
		select {
		case <-stop.Main:
			close(stop.Slave)
			return
		case <-in.Slave:
			out.WaitWrite()
			in.ConfirmRead()
			out.ConfirmWrite()
			out.Write(in.Read())
		}
	}
}

func SyncdrainChannel(in1, in2, stop Port) {
	for {
		select {
		case <-stop.Main:
			close(stop.Slave)
			return
		case <-in1.Slave:
			select {
			case <-stop.Main:
				close(stop.Slave)
				return
			case <-in2.Slave:
			}
		}
		in2.WaitRead()
		in1.ConfirmRead()
		in2.ConfirmRead()
		in1.Read()
		in2.Read()
	}
}

func FifoChannel(in, out, stop Port) {
	for {
		select {
		case <-stop.Main:
			close(stop.Slave)
			return
		case <-in.Slave:
			in.ConfirmRead()
			c := in.Read()
			out.SyncWrite(c)
		}
	}
}

func LossysyncChannel(in, out, stop Port) {
	for {
		select {
		case <-stop.Main:
			close(stop.Slave)
			return
		}
		// FIXME the SyncRead operation may blocks this channel
		// and hence it cannot be closed by stop Port
		c := in.SyncRead()
		select {
		// try WaitWrite
		case out.Slave <- "write":
			out.ConfirmWrite()
			out.Write(c)
		default:
			// do nothing
		}
	}
}

func MergerChannel(in1, in2, out, stop Port) {
	for {
		// considering the syntax of select, here we use
		// <-in.slave instead of in.WaitRead()
		select {
		case <-stop.Main:
			close(stop.Slave)
			return
		case <-in1.Slave:
			out.WaitWrite()
			in1.ConfirmRead()
			out.ConfirmWrite()
			out.Write(in1.Read())
		case <-in2.Slave:
			out.WaitWrite()
			in2.ConfirmRead()
			out.ConfirmWrite()
			out.Write(in2.Read())
		}
	}
}

func ReplicatorChannel(in Port, out Ports, stop Port) {
	for {
		select {
		case <-stop.Main:
			close(stop.Slave)
			return
		case <-in.Slave:
			out.WaitWrite()
			in.ConfirmRead()
			out.ConfirmWrite()
			out.Write(in.Read())
		}
	}
}

func BufferChannel(in, out, stop Port) {
	buf := []string{}
	for {
		select {
		case <-stop.Main:
			close(stop.Slave)
			return
		case <-in.Slave:
			in.ConfirmRead()
			buf = append(buf, in.Read())
		case out.Slave <- "write":
			out.ConfirmWrite()
			out.Write(buf[0])
			buf = buf[1:]
		}
	}
}
