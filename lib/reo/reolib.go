package reo

func SyncChannel(in, out, stop Port) {
	for {
		select {
		case <-stop.Main:
			close(stop.Slave)
			return
		default:
		}
		in.WaitRead()
		out.WaitWrite()
		in.ConfirmRead()
		out.ConfirmWrite()
		out.Write(in.Read())
	}
}

func SyncdrainChannel(in1, in2, stop Port) {
	for {
		select {
		case <-stop.Main:
			close(stop.Slave)
			return
		default:
		}
		in1.WaitRead()
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
		default:
		}
		c := in.SyncRead()
		out.SyncWrite(c)
	}
}

func LossysyncChannel(in, out, stop Port) {
	for {
		select {
		case <-stop.Main:
			close(stop.Slave)
			return
		}
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
		select {
		case <-stop.Main:
			close(stop.Slave)
			return
		default:
		}
		// considering the syntax of select, here we use
		// <-in.slave instead of in.WaitRead()
		select {
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

func ReplicatorChannel(in, stop Port, out Ports) {
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
